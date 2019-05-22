package driver

import (
	"context"
	"os"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubevault/csi-driver/pkg/util"
	vs "github.com/kubevault/operator/pkg/vault/secret"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InstanceNotFound = errors.New("instance not found")
	SecretEngineKey  = "engine"
)

var (
	nodeVolumeTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "csi_node_volume_total",
		Help: "Total number of volume at current stage",
	}, []string{"volume_name", "type"})

	vaultFetchSecretTotal = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "csi_fetch_secret_total",
		Help: "Total number of fetched secret from vault driver",
	}, []string{"volume_name", "type"})

	vaultRenewDurations = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "vault_renew_secret_duration_seconds",
		Help: "The time duration of renewing secret",
	})

	vaultFetchSecretDurations = prometheus.NewSummary(prometheus.SummaryOpts{
		Name: "vault_fetch_secret_duration_seconds",
		Help: "The time duration of fetching secret",
	})
)

func init() {
	prometheus.MustRegister(nodeVolumeTotal, vaultFetchSecretTotal, vaultRenewDurations, vaultFetchSecretDurations)
}

// NodeStageVolume mounts the volume to a staging path on the node. This is
// called by the CO before NodePublishVolume and is used to temporary mount the
// volume to a staging path. Once mounted, NodePublishVolume will make sure to
// mount it to the appropriate path
func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeStageVolume Volume ID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeStageVolume Staging Target Path must be provided")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "NodeStageVolume Volume Capability must be provided")
	}

	d.log.WithFields(logrus.Fields{
		"request": req,
		"method":  "node_stage_volume",
	}).Info("node stage volume called")

	options := req.VolumeContext

	if _, ok := options[SecretEngineKey]; !ok {
		return nil, errors.Errorf("Missing engine name field (%s) in options %v", SecretEngineKey, options)
	}

	_, sKey := options[vs.SecretKey]
	_, rKey := options[vs.RoleKey]
	if !sKey && !rKey {
		return nil, errors.Errorf("Misssing secret field (%s/%s)", vs.SecretKey, vs.RoleKey)
	}

	if _, ok := options[vs.PathKey]; !ok {
		return nil, errors.Errorf("Misssing secret name field (%s)", vs.PathKey)
	}

	fsType := "tmpfs"

	ll := d.log.WithFields(logrus.Fields{
		"volume_id":           req.VolumeId,
		"staging_target_path": req.StagingTargetPath,
		"fsType":              fsType,
		"mount_options":       options,
		"method":              "node_stage_volume",
	})

	formatted, err := d.mounter.IsFormatted(req.StagingTargetPath)
	if err != nil {
		return nil, err
	}

	if !formatted {
		ll.Info("formatting the volume for staging")
		if err := d.mounter.Format(req.StagingTargetPath, fsType); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		ll.Info("source device is already formatted")
	}
	if err := os.MkdirAll(req.StagingTargetPath, 0755); err != nil {
		return nil, err
	}

	ll.Info("formatting and mounting stage volume is finished")
	nodeVolumeTotal.WithLabelValues(req.VolumeId, "stage").Inc()
	return &csi.NodeStageVolumeResponse{}, nil
}

// NodeUnstageVolume unstages the volume from the staging path
func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeUnstageVolume Volume ID must be provided")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeUnstageVolume Staging Target Path must be provided")
	}
	ll := d.log.WithFields(logrus.Fields{
		"request": req,
		"method":  "node_unstage_volume",
	})
	ll.Info("node unstage volume called")
	resp := &csi.NodeUnstageVolumeResponse{}

	mounted, err := d.mounter.IsMounted("", req.StagingTargetPath)
	if err != nil {
		return nil, err
	}

	if mounted {
		ll.Info("unmounting the staging target path")
		err := d.mounter.VaultUnmount(req.StagingTargetPath)
		if err != nil {
			return resp, err
		}

	} else {
		ll.Info("staging target path is already unmounted")
	}
	nodeVolumeTotal.WithLabelValues(req.VolumeId, "unstage").Inc()
	return &csi.NodeUnstageVolumeResponse{}, err
}

// NodePublishVolume mounts the volume mounted to the staging path to the target path
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"request": req,
		"method":  "node_publish_volume",
	}).Info("node publish volume called")
	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodePublishVolume Staging Target Path must be provided")
	}

	if req.TargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "NodePublishVolume Target Path must be provided")
	}

	if len(req.VolumeContext) == 0 {
		return nil, status.Error(codes.InvalidArgument, "NodePublishVolume Volume attributes are not provided")
	}

	podInfo, err := getPodInfo(req.VolumeContext)
	if err != nil {
		return nil, err
	}

	podInfo.KubeClient = d.kubeClient
	podInfo.AppClient = d.appClient

	podInfo.RefNamespace, podInfo.RefName, err = getAppBindingInfo(req.VolumeContext)
	if err != nil {
		return nil, err
	}

	source := req.StagingTargetPath
	target := req.TargetPath
	fsType := "tmpfs"
	//mnt := req.VolumeCapability.GetMount()

	ll := d.log.WithFields(logrus.Fields{
		"volume_id": req.VolumeId,
		"source":    source,
		"target":    target,
		"method":    "node_publish_volume",
	})
	opts := []string{"rw"}
	mounted, err := d.mounter.IsMounted(source, target)
	if err != nil {
		return nil, err
	}
	if !mounted {
		ll.Info("mounting the volume")
		if err := d.mounter.Mount("tmpfs", target, fsType, opts...); err != nil {
			return nil, status.Error(codes.Internal, err.Error())
		}
	} else {
		ll.Info("volume is already mounted")
	}

	options := req.VolumeContext

	// login with policy token

	authClient, err := util.NewVaultClient(podInfo)
	if err != nil {
		return nil, err
	}

	engineName, ok := options[SecretEngineKey]
	if !ok {
		return nil, errors.Errorf("Empty engine name")
	}

	now := time.Now()
	secretData, err := util.FetchSecret(engineName, authClient, options, target)
	if err != nil {
		return nil, err
	}
	vaultFetchSecretDurations.Observe(float64(time.Since(now).Nanoseconds() / int64(time.Microsecond)))
	vaultFetchSecretTotal.WithLabelValues(req.VolumeId, "Fetch").Inc()

	ll.Info("mounting the volume with ro")
	if err := d.mounter.Mount("tmpfs", target, fsType, []string{"remount,ro"}...); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if _, found := d.ch[req.VolumeId]; !found {
		if secretData.Renewable {
			vaultFetchSecretTotal.WithLabelValues(req.VolumeId, "Renew").Inc()
			start := time.Now()
			renewer, err := util.SetRenewal(authClient, secretData)
			if err != nil {
				return nil, err
			}
			d.ch[req.VolumeId] = renewer

			vaultRenewDurations.Observe(float64(time.Since(start).Nanoseconds() / int64(time.Microsecond)))
		}
	}
	nodeVolumeTotal.WithLabelValues(req.VolumeId, "publish").Inc()

	ll.Info("bind mounting the volume is finished")
	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume unmounts the volume from the target path
func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	d.log.WithFields(logrus.Fields{
		"request": req,
		"method":  "node_unpublish_volume",
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
		if _, err := os.Stat(req.TargetPath); err == nil {
			os.RemoveAll(req.TargetPath)
		}
	} else {
		ll.Info("target path is already unmounted")
	}

	if _, found := d.ch[req.VolumeId]; found {
		vaultFetchSecretTotal.WithLabelValues(req.VolumeId, "Renew").Desc()
		util.StopRenew(d.ch[req.VolumeId])
	}

	nodeVolumeTotal.WithLabelValues(req.VolumeId, "unpublish").Inc()
	ll.Info("unmounting volume is finished")
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

func (d *Driver) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	d.log.WithField("method", "node_get_id").Info("node get id called")
	return nil, status.Error(codes.Unimplemented, "NodeGetVolumeStats is not implemented")
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
		NodeId:            d.NodeId,
		MaxVolumesPerNode: 10,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				driverName: d.NodeId,
			},
		},
	}, nil
}

func (d *Driver) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	d.log.WithField("method", "node_expand_volume").Info("node expand volume called")
	return nil, status.Error(codes.Unimplemented, "NodeExpandVolume is not implemented")
}
