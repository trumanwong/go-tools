package volcenginecloud

import (
	"context"
	"github.com/volcengine/ve-tos-golang-sdk/v2/tos"
)

type OssClient struct {
	client *tos.ClientV2
}

func NewOssClient(accessKey, secretKey, endpoint, region string) (*OssClient, error) {
	client, err := tos.NewClientV2(
		endpoint,
		tos.WithCredentials(tos.NewStaticCredentials(accessKey, secretKey)),
		tos.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}
	return &OssClient{client: client}, nil
}

// HeadObject 判断对象是否存在
func (c OssClient) HeadObject(ctx context.Context, req *tos.HeadObjectV2Input) (*tos.HeadObjectV2Output, error) {
	return c.client.HeadObjectV2(ctx, req)
}

// DeleteObject 删除对象
func (c OssClient) DeleteObject(ctx context.Context, req *tos.DeleteObjectV2Input) (*tos.DeleteObjectV2Output, error) {
	return c.client.DeleteObjectV2(ctx, req)
}

// RenameObject 重命名对象
func (c OssClient) RenameObject(ctx context.Context, req *tos.RenameObjectInput) (*tos.RenameObjectOutput, error) {
	return c.client.RenameObject(ctx, req)
}

// RestoreObject 恢复对象
func (c OssClient) RestoreObject(ctx context.Context, req *tos.RestoreObjectInput) (*tos.RestoreObjectOutput, error) {
	return c.client.RestoreObject(ctx, req)
}

// ListObjectsType2 列举对象
func (c OssClient) ListObjectsType2(ctx context.Context, req *tos.ListObjectsType2Input) (*tos.ListObjectsType2Output, error) {
	return c.client.ListObjectsType2(ctx, req)
}

// PutObject 普通上传
func (c OssClient) PutObject(ctx context.Context, req *tos.PutObjectV2Input) (*tos.PutObjectV2Output, error) {
	return c.client.PutObjectV2(ctx, req)
}

// AppendObject 追加上传
func (c OssClient) AppendObject(ctx context.Context, req *tos.AppendObjectV2Input) (*tos.AppendObjectV2Output, error) {
	return c.client.AppendObjectV2(ctx, req)
}

// PreSignedPostSignature Post表单预签名
func (c OssClient) PreSignedPostSignature(ctx context.Context, req *tos.PreSingedPostSignatureInput) (*tos.PreSingedPostSignatureOutput, error) {
	return c.client.PreSignedPostSignature(ctx, req)
}
