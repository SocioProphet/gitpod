// Copyright (c) 2021 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License-AGPL.txt in the project root for license information.

package service

import (
	"context"
	"fmt"

	"github.com/gitpod-io/gitpod/common-go/log"
	"github.com/gitpod-io/gitpod/content-service/api"
	"github.com/gitpod-io/gitpod/content-service/pkg/storage"
)

// ContentService implements ContentServiceServer
type ContentService struct {
	cfg storage.Config
}

// NewContentService create a new content service
func NewContentService(cfg storage.Config) (res *ContentService, err error) {
	return &ContentService{cfg}, nil
}

// DownloadUrl provides a upload URL
func (cs *ContentService) UploadUrl(ctx context.Context, req *api.UploadUrlRequest) (resp *api.UploadUrlResponse, err error) {
	log.Info("UploadUrl")
	return &api.UploadUrlResponse{
		Url:         req.Ref.Name,
		CommitToken: "123",
	}, nil
}

// DownloadUrl provides a download URL
func (cs *ContentService) DownloadUrl(ctx context.Context, req *api.DownloadUrlRequest) (resp *api.DownloadUrlResponse, err error) {
	log.Infof("DownloadUrl %s %s", req.Ref.OwnerId, req.Ref.Name)

	s, err := storage.NewPresignedAccess(&cs.cfg)
	if err != nil {
		return nil, err
	}

	bucket := fmt.Sprintf("gitpod-user-%s", req.Ref.OwnerId) // FIXME
	object := req.Ref.Name                                   // TODO: remove me
	// object := fmt.Sprintf("blobs/%s", req.Ref.Name)          // TODO: validate name
	info, err := s.SignDownload(ctx, bucket, object)
	if err != nil {
		return nil, err
	}

	return &api.DownloadUrlResponse{
		Url: info.URL,
	}, nil
}

// Delete deletes a blob
func (s *ContentService) Delete(ctx context.Context, req *api.DeleteRequest) (resp *api.DeleteResponse, err error) {
	log.Info("Delete")
	return &api.DeleteResponse{}, nil
}
