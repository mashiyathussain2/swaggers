package app

import (
	"go-app/mock"
	"go-app/model"
	"go-app/schema"
	"go-app/server/kafka"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/rs/zerolog"
	segKafka "github.com/segmentio/kafka-go"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"syreclabs.com/go/faker"
)

func TestContentUpdateProcessor_ProcessBrandMessage(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		Logger *zerolog.Logger
	}
	type args struct {
		msg kafka.Message
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		prepare    func(tt *TC)
		buildStubs func(tt *TC, m *mock.MockContent)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, m *mock.MockContent) {
				id, _ := primitive.ObjectIDFromHex("6052e43c29fc71ce32cf7772")
				opts := &schema.UpdateContentBrandInfoOpts{
					ID:   id,
					Name: "test 2 2",
					Logo: &model.IMG{
						Height: 225,
						Width:  225,
						SRC:    "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRMOBZQgE7t6Dr_Ompn0EaPk1pShjRWUA_dIA\u0026usqp=CAU",
					},
				}
				m.EXPECT().UpdateContentBrandInfo(opts)
			},
			prepare: func(tt *TC) {
				value := string(`{"meta":{"_id":{"$oid":"6052e43c29fc71ce32cf7772"},"ts":{"$timestamp":{"t":1616055057,"i":2}},"ns":"entity.brand","op":"u","updates":{"removed":[],"changed":{"updated_at":{"$date":"2021-03-18T08:10:57.432Z"}}}},"data":{"website":"https://test.com","updated_at":{"$date":"2021-03-18T08:10:57.432Z"},"domain":"vasu","cover_img":{"width":225,"src":"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcTK2No4UYBzFd5CBUz-ZPFp1Pfh66Sh4RIwUg\u0026usqp=CAU","height":225},"_id":{"$oid":"6052e43c29fc71ce32cf7772"},"name":"test 2 2","logo":{"height":225,"width":225,"src":"https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRMOBZQgE7t6Dr_Ompn0EaPk1pShjRWUA_dIA\u0026usqp=CAU"},"created_at":{"$date":"2021-03-18T05:25:15.791Z"},"lname":"test 2 2","registered_name":"test pvt ltd","fulfillment_email":"test@fulfillment.com"}}`)
				tt.args.msg = segKafka.Message{
					Topic:     "entity.brand",
					Partition: 0,
					Offset:    faker.RandomInt64(1, 100),
					Key:       []byte(primitive.NewObjectID().Hex()),
					Value:     []byte(value),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContent := mock.NewMockContent(ctrl)
			csp := &ContentUpdateProcessor{
				App:    tt.fields.App,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = mockContent
			tt.fields.App.ContentUpdateProcessor = csp
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockContent)
			csp.ProcessBrandMessage(tt.args.msg)
		})
	}
}

func TestContentUpdateProcessor_ProcessContentMessage(t *testing.T) {
	t.Parallel()

	app := NewTestApp(getTestConfig())
	defer CleanTestApp(app)

	type fields struct {
		App    *App
		Logger *zerolog.Logger
	}
	type args struct {
		msg kafka.Message
	}
	type TC struct {
		name       string
		fields     fields
		args       args
		prepare    func(tt *TC)
		buildStubs func(*TC, *mock.MockContent, *mock.MockProducer, *mock.MockMedia)
	}
	tests := []TC{
		{
			name: "[Ok]",
			fields: fields{
				App:    app,
				Logger: app.Logger,
			},
			args: args{},
			buildStubs: func(tt *TC, m *mock.MockContent, m2 *mock.MockProducer, m3 *mock.MockMedia) {
				Id, _ := primitive.ObjectIDFromHex("605337b329fc71ce32cf7773")
				resp := []model.BrandInfo{
					{
						ID:   Id,
						Name: "Hypd Store",
						Logo: &model.IMG{
							SRC:    "http://robohash.org/ovr2inpgLMQNC3Bh.png?size=100x100",
							Width:  50,
							Height: 50,
						},
					},
				}

				mediaId, _ := primitive.ObjectIDFromHex("605337e8b0925ee5b5755285")
				t1, _ := time.Parse(time.RFC3339, "2021-02-18T00:00:00+00:00")
				resp1 := schema.GetMediaResp{
					ID:       mediaId,
					FileName: "testfile.png",
					Dimensions: &model.Dimensions{
						Width:  100,
						Height: 100,
					},
					URL:       "http://robohash.org/ovr2inpgLMQNC3Bh.png?size=100x100",
					CreatedAt: t1,
				}
				catalogID, _ := primitive.ObjectIDFromHex("604cf084dd7faf4bfe9c6d01")
				resp2 := []model.CatalogInfo{
					{
						ID:   catalogID,
						Name: "test Product",
						FeaturedImage: &model.IMG{
							SRC:    "http://robohash.org/yOF8RpCUfN1Vodrc.png?size=100x100",
							Width:  100,
							Height: 100,
						},
						BasePrice:   model.SetINRPrice(100),
						RetailPrice: model.SetINRPrice(100),
					},
				}

				m3.EXPECT().GetImageMediaByID(mediaId).Times(1).Return(&resp1, nil)
				m.EXPECT().GetBrandInfo([]string{"605337b329fc71ce32cf7773"}).Times(1).Return(resp, nil)
				m.EXPECT().GetCatalogInfo([]string{"604cf084dd7faf4bfe9c6d01"}).Times(1).Return(resp2, nil)

				kafkaMessage := segKafka.Message{
					Key:   []byte("605341a5a868e897c05ea6de"),
					Value: []byte(`{"_id":"605341a5a868e897c05ea6de","type":"catalog_content","media_type":"image","media_id":"605337e8b0925ee5b5755285","media_info":{"id":"605337e8b0925ee5b5755285","filename":"testfile.png","created_at":"2021-02-18T00:00:00Z","dimensions":{"height":100,"width":100},"url":"http://robohash.org/ovr2inpgLMQNC3Bh.png?size=100x100"},"brand_ids":["605337b329fc71ce32cf7773"],"brand_info":[{"id":"605337b329fc71ce32cf7773","name":"Hypd Store","logo":{"src":"http://robohash.org/ovr2inpgLMQNC3Bh.png?size=100x100","height":50,"width":50}}],"label":{"interests":["test"],"genders":["M","F"]},"is_processed":true,"is_active":false,"catalog_ids":["604cf084dd7faf4bfe9c6d01"],"catalog_info":[{"id":"604cf084dd7faf4bfe9c6d01","name":"test Product","featured_image":{"src":"http://robohash.org/yOF8RpCUfN1Vodrc.png?size=100x100","height":100,"width":100},"base_price":{"iso":"inr","value":100},"retail_price":{"iso":"inr","value":100}}],"created_at":"2021-03-18T17:33:49.501+05:30","processed_at":"0001-01-01T00:00:00Z","updated_at":"0001-01-01T00:00:00Z"}`),
				}
				m2.EXPECT().Publish(kafkaMessage)
			},
			prepare: func(tt *TC) {
				value := string(`{
					"meta": {
					  "_id": {
						"$oid": "605341a5a868e897c05ea6de"
					  },
					  "ts": {
						"$timestamp": {
						  "t": 1616069030,
						  "i": 1
						}
					  },
					  "ns": "cms.content",
					  "op": "i"
					},
					"data": {
					  "media_id": {
						"$oid": "605337e8b0925ee5b5755285"
					  },
					  "is_active": false,
					  "type": "catalog_content",
					  "media_type": "image",
					  "is_processed": true,
					  "created_at": {
						"$date": "2021-03-18T12:03:49.501Z"
					  },
					  "_id": {
						"$oid": "605341a5a868e897c05ea6de"
					  },
					  "label": {
						"interests": [
						  "test"
						],
						"genders": [
						  "M",
						  "F"
						]
					  },
					  "catalog_ids": [
						{
						  "$oid": "604cf084dd7faf4bfe9c6d01"
						}
					  ],
					  "brand_ids": [
						{
						  "$oid": "605337b329fc71ce32cf7773"
						}
					  ]
					}
				  }`)
				tt.args.msg = segKafka.Message{
					Topic:     "cms.content",
					Partition: 0,
					Offset:    faker.RandomInt64(1, 100),
					Key:       []byte(primitive.NewObjectID().Hex()),
					Value:     []byte(value),
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mockContent := mock.NewMockContent(ctrl)
			mockProducer := mock.NewMockProducer(ctrl)
			mockMedia := mock.NewMockMedia(ctrl)
			csp := &ContentUpdateProcessor{
				App:    tt.fields.App,
				Logger: tt.fields.Logger,
			}
			tt.fields.App.Content = mockContent
			tt.fields.App.Media = mockMedia
			tt.fields.App.ContentFullProducer = mockProducer
			tt.fields.App.ContentUpdateProcessor = csp
			tt.prepare(&tt)
			tt.buildStubs(&tt, mockContent, mockProducer, mockMedia)
			csp.ProcessContentMessage(tt.args.msg)
		})
	}
}
