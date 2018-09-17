package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi/v0"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	InstanceNotFound = errors.New("instance not found")
)

// NodeStageVolume mounts the volume to a staging path on the node. This is
// called by the CO before NodePublishVolume and is used to temporary mount the
// volume to a staging path. Once mounted, NodePublishVolume will make sure to
// mount it to the appropriate path
func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	d.log.Info("node stage volume called")
	mnt := req.VolumeCapability.GetMount()
	options := mnt.MountFlags

	fsType := "tmpfs"
	if mnt.FsType != "" {
		fsType = mnt.FsType
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

	/*
		formatted, err := d.mounter.IsFormatted(source)
		if err != nil {
			return nil, err
		}

		if !formatted {
			ll.Info("formatting the volume for staging")
			if err := d.mounter.Format(source, fsType); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			ll.Info("source device is already formatted")
		}

		ll.Info("mounting the volume for staging")

		mounted, err := d.mounter.IsMounted(source, target)
		if err != nil {
			return nil, err
		}

		if !mounted {
			if err := d.mounter.Mount(source, target, fsType, options...); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
		} else {
			ll.Info("source device is already mounted to the target path")
		}
	*/
	ll.Info("formatting and mounting stage volume is finished")
	return &csi.NodeStageVolumeResponse{}, nil
	return nil, nil
}

// NodeUnstageVolume unstages the volume from the staging path
func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"method": "node_unstage_volume",
	}).Info("node unstage volume called")
	return nil, nil
}

// NodePublishVolume mounts the volume mounted to the staging path to the target path
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"method": "node_publish_volume",
	}).Info("node publish volume called")
	return nil, nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"method": "node_unpublish_volume",
	}).Info("node unpublish volume called")
	return nil, nil
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
		MaxVolumesPerNode: 5,
	}, nil
}
