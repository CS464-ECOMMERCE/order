package models

import (
	"github.com/lib/pq"
)

type Product struct {
	Id              uint64         `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Name            string         `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Price           float32        `protobuf:"fixed32,3,opt,name=price,proto3" json:"price,omitempty"`
	Inventory       uint64         `protobuf:"varint,4,opt,name=inventory,proto3" json:"inventory,omitempty"`
	Description     string         `protobuf:"bytes,5,opt,name=description,proto3" json:"description,omitempty"`
	Images          pq.StringArray `protobuf:"bytes,6,rep,name=images,proto3" json:"images,omitempty" gorm:"type:text[]"` // [bucketname]
	StripePriceId   string         `protobuf:"bytes,7,opt,name=stripe_price_id,json=stripePriceId,proto3" json:"stripe_price_id,omitempty"`
	StripeProductId string         `protobuf:"bytes,8,opt,name=stripe_product_id,json=stripeProductId,proto3" json:"stripe_product_id,omitempty"`
	MerchantId      uint64         `protobuf:"varint,9,opt,name=merchant_id,json=merchantId,proto3" json:"merchant_id,omitempty"`
}
