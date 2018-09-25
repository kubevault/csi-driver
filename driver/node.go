package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"fmt"
	"strings"
	"github.com/kubevault/csi-driver/vault"
)

var (
	InstanceNotFound = errors.New("instance not found")
)

// NodeStageVolume mounts the volume to a staging path on the node. This is
// called by the CO before NodePublishVolume and is used to temporary mount the
// volume to a staging path. Once mounted, NodePublishVolume will make sure to
// mount it to the appropriate path
func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"request": req,
		"method": "node_stage_volume",
	}).Info("node stage volume called")
	//mnt := req.VolumeCapability.GetMount()
	//options := mnt.MountFlags
	options := req.VolumeAttributes

	fsType := "tmpfs"
	if v, ok := options["fsType"]; ok {
		fsType = v
	}

	ll := d.log.WithFields(logrus.Fields{
		"volume_id": req.VolumeId,
		//"volume_name":         vol.Label,
		"staging_target_path": req.StagingTargetPath,
		//	"source":              source,
		"fsType":        fsType,
		"mount_options": options,
		"method":        "node_stage_volume",
	})

	stringPolicies, ok := options["policy"]
	if !ok {
		return nil, errors.Errorf("Missing policies")
	}
	policies := strings.Split(strings.Replace(stringPolicies, " ", "", -1), ",")
	if len(policies) == 0 {
		return nil, errors.Errorf("Empty policies")
	}
	token, err := d.vaultClient.GetPolicyToken(policies, true)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	// login with policy token

	client, err := vault.NewVaultClient(d.url, token, nil)

	fmt.Println(client)

	if err := d.mounter.VaultMount(req.StagingTargetPath, fsType, options); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	ll.Info("formatting and mounting stage volume is finished")
	return &csi.NodeStageVolumeResponse{}, nil
	return nil, nil
}

// NodeUnstageVolume unstages the volume from the staging path
func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"request": req,
		"method": "node_unstage_volume",
	}).Info("node unstage volume called")
	err:= d.mounter.VaultUnmount(req.StagingTargetPath)
	fmt.Println(err)

	return &csi.NodeUnstageVolumeResponse{}, nil
}

// NodePublishVolume mounts the volume mounted to the staging path to the target path
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"request": req,
		"method": "node_publish_volume",
	}).Info("node publish volume called")
	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodePublishVolume Staging Target Path must be provided")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodePublishVolume Target Path must be provided")
	}

	source := req.StagingTargetPath
	target := req.TargetPath

	mnt := req.VolumeCapability.GetMount()
	options := mnt.MountFlags

	// TODO(arslan): linode we need bind here? check it out
	// Perform a bind mount to the full path to allow duplicate mounts of the same PD.
	options = append(options, "bind")
	if req.Readonly {
		options = append(options, "ro")
	}

	fsType := "tmpfs"
	if mnt.FsType != "" {
		fsType = mnt.FsType
	}

	ll := d.log.WithFields(logrus.Fields{
		"volume_id":     req.VolumeId,
		"source":        source,
		"target":        target,
		"fsType":        fsType,
		"mount_options": options,
		"method":        "node_publish_volume",
	})


	mounted, err := d.mounter.IsMounted(source, target)
	if err != nil {
		return nil, err
	}

	if !mounted {
		ll.Info("mounting the volume")
		if err := d.mounter.Mount(source, target, fsType, options...); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		ll.Info("volume is already mounedt")
	}

	ll.Info("bind mounting the volume is finished")
	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"request": req,
		"method": "node_unpublish_volume",
	}).Info("node unpublish volume called")
	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeUnpublishVolume Target Path must be provided")
	}

	ll := d.log.WithFields(logrus.Fields{
		"volume_id":   req.VolumeId,
		"target_path": req.TargetPath,
		"method":      "node_unpublish_volume",
	})
	ll.Info("node unpublish volume called")

	mounted, err := d.mounter.IsMounted("", req.TargetPath)
	if err != nil {
		return nil, err
	}

	if mounted {
		ll.Info("unmounting the target path")
		err := d.mounter.Unmount(req.TargetPath)
		if err != nil {
			return nil, err
		}
	} else {
		ll.Info("target path is already unmounted")
	}

	ll.Info("unmounting volume is finished")
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetId returns the unique id of the node. This should eventually return
// the linode ID if possible. This is used so the CO knows where to place the
// workload. The result of this function will be used by the CO in
// ControllerPublishVolume.
func (d *Driver) NodeGetId(ctx context.Context, req *csi.NodeGetIdRequest) (*csi.NodeGetIdResponse, error) {
	d.log.WithField("method", "node_get_id").Info("node get id called")
	return &csi.NodeGetIdResponse{
		NodeId: d.nodeId,
	}, nil
}

// NodeGetCapabilities returns the supported capabilities of the node server
func (d *Driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	// currently there is a single NodeServer capability according to the spec
	nscap := &csi.NodeServiceCapability{
		Type: &csi.NodeServiceCapability_Rpc{
			Rpc: &csi.NodeServiceCapability_RPC{
				Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
			},
		},
	}

	d.log.WithFields(logrus.Fields{
		"node_capabilities": nscap,
		"method":            "node_get_capabilities",
	}).Info("node get capabilities called")
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			nscap,
		},
	}, nil
}

func (d *Driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	d.log.WithField("method", "node_get_info").Info("node get info called")
	return &csi.NodeGetInfoResponse{
		NodeId:            d.nodeId,
		MaxVolumesPerNode: 10,
	}, nil
}
