// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package service

import (
	"context"
	"strings"

	"github.com/gitpod-io/gitpod/common-go/log"
	"github.com/gitpod-io/gitpod/content-service/api"
	"github.com/gitpod-io/gitpod/content-service/pkg/storage"
)

// ContentService implements ContentServiceServer
type ContentService struct {
	cfg storage.Config
	s   storage.PresignedAccess
}

// NewContentService create a new content service
func NewContentService(cfg storage.Config) (res *ContentService, err error) {
	s, err := storage.NewPresignedAccess(&cfg)
	if err != nil {
		return nil, err
	}
	return &ContentService{cfg, s}, nil
}

// UploadUrl provides a upload URL
func (cs *ContentService) UploadUrl(ctx context.Context, req *api.UploadUrlRequest) (resp *api.UploadUrlResponse, err error) {
	log.Info("UploadUrl") // TODO: trace

	blobName, err := cs.s.BlobObject(req.Ref.Name)
	if err != nil {
		return nil, err
	}

	info, err := cs.s.SignUpload(ctx, cs.s.Bucket(req.Ref.OwnerId), blobName)
	if err != nil {
		return nil, err
	}

	return &api.UploadUrlResponse{
		Url: cs.postprocessUrl(info.URL),
	}, nil
}

// DownloadUrl provides a download URL
func (cs *ContentService) DownloadUrl(ctx context.Context, req *api.DownloadUrlRequest) (resp *api.DownloadUrlResponse, err error) {
	log.Infof("DownloadUrl %s %s", req.Ref.OwnerId, req.Ref.Name) // TODO: trace

	blobName, err := cs.s.BlobObject(req.Ref.Name)
	if err != nil {
		return nil, err
	}

	info, err := cs.s.SignDownload(ctx, cs.s.Bucket(req.Ref.OwnerId), blobName)
	if err != nil {
		return nil, err
	}

	return &api.DownloadUrlResponse{
		Url: cs.postprocessUrl(info.URL),
	}, nil
}

// Delete deletes a blob
func (cs *ContentService) Delete(ctx context.Context, req *api.DeleteRequest) (resp *api.DeleteResponse, err error) {
	log.Info("Delete") // TODO: trace

	blobName, err := cs.s.BlobObject(req.Ref.Name)
	if err != nil {
		return nil, err
	}

	err = cs.s.DeleteObject(ctx, cs.s.Bucket(req.Ref.OwnerId), blobName)
	if err != nil {
		return nil, err
	}

	return &api.DeleteResponse{}, nil
}

// TODO: is there a better concept?
func (cs *ContentService) postprocessUrl(url string) string {
	if cs.cfg.ReplaceHost != "" && cs.cfg.Kind == storage.MinIOStorage {
		return strings.Replace(url, cs.cfg.MinIOConfig.Endpoint, cs.cfg.ReplaceHost+"/minio", 1)
	}
	return url
}
