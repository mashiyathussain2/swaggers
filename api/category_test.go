package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-app/mock"
	"go-app/schema"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"syreclabs.com/go/faker"
)

func TestAPI_createCategory(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	opts := schema.GetRandomCreateCategoryOpts()
	resp := schema.GetRandomCreateCategoryResp(opts)
	validBody, _ := json.Marshal(opts)

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Error] No Body",
			url:    "/api/keeper/category",
			method: http.MethodPost,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().CreateCategory(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, string(data), string("{\"error\":[{\"message\":\"Request body must not be empty\",\"type\":\"BadRequest\"}],\"success\":false,\"request_id\":\"\"}\n"))
			},
		},
		{
			name:   "[Ok]",
			url:    "/api/keeper/category",
			method: http.MethodPost,
			body:   bytes.NewReader(validBody),
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().CreateCategory(gomock.Any()).Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_editCategory(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	opts := schema.GetRandomEditCategoryOpts()
	resp := schema.GetRandomEditCategoryResp(opts)
	validBody, _ := json.Marshal(opts)

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Error] No Body",
			url:    "/api/keeper/category",
			method: http.MethodPut,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().CreateCategory(gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusBadRequest, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				assert.JSONEq(t, string(data), string("{\"error\":[{\"message\":\"Request body must not be empty\",\"type\":\"BadRequest\"}],\"success\":false,\"request_id\":\"\"}\n"))
			},
		},
		{
			name:   "[Ok]",
			url:    "/api/keeper/category",
			method: http.MethodPut,
			body:   bytes.NewReader(validBody),
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().EditCategory(gomock.Any()).Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusCreated, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_getCategory(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	var resp []schema.GetCategoriesResp
	for i := 0; i < faker.RandomInt(0, 5); i++ {
		resp = append(resp, *schema.GetRandomGetCategoriesResp())
	}

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Ok]",
			url:    "/api/keeper/category",
			method: http.MethodGet,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().GetCategories().Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_getCategoryMain(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	resp := make(map[string]schema.GetMainCategoriesMapResp)
	for i := 0; i < faker.RandomInt(0, 5); i++ {
		obj := schema.GetRandomGetMainCategoriesMapResp()
		resp[obj.ID.Hex()] = *obj
	}

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Ok]",
			url:    "/api/keeper/category/main",
			method: http.MethodGet,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().GetMainCategoriesMap().Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_getParentCategory(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	var resp []schema.GetParentCategoriesResp
	for i := 0; i < faker.RandomInt(0, 5); i++ {
		resp = append(resp, *schema.GetRandomGetParentCategoriesResp())
	}

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Ok]",
			url:    "/api/category/lvl1",
			method: http.MethodGet,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().GetMainParentCategories().Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_getMainCategoryByParentID(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	var resp []schema.GetMainCategoriesByParentIDResp
	id := primitive.NewObjectIDFromTimestamp(time.Now())
	for i := 0; i < faker.RandomInt(0, 5); i++ {
		resp = append(resp, *schema.GetRandomGetMainCategoriesByParentIDResp())
	}

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Ok]",
			url:    fmt.Sprintf("/api/category/%s/lvl2", id.Hex()),
			method: http.MethodGet,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().GetMainCategoriesByParentID(id).Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}

func TestAPI_getSubCatergoryByParentID(t *testing.T) {

	api := NewTestAPI(getTestConfig())

	var resp []schema.GetSubCategoriesByParentIDResp
	id := primitive.NewObjectIDFromTimestamp(time.Now())

	for i := 0; i < faker.RandomInt(0, 5); i++ {
		resp = append(resp, *schema.GetRandomGetSubCategoriesByParentIDResp())
	}

	tests := []struct {
		name          string
		url           string
		method        string
		body          io.Reader
		buildStubs    func(ex *mock.MockCategory)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:   "[Ok]",
			url:    fmt.Sprintf("/api/category/%s/lvl3", id.Hex()),
			method: http.MethodGet,
			body:   nil,
			buildStubs: func(m *mock.MockCategory) {
				m.EXPECT().GetSubCategoriesByParentID(id).Times(1).Return(resp, nil)
			},
			checkResponse: func(t *testing.T, r *httptest.ResponseRecorder) {
				assert.Equal(t, http.StatusOK, r.Code)
				data, err := ioutil.ReadAll(r.Body)
				assert.Nil(t, err)
				respBody, _ := json.Marshal(map[string]interface{}{"success": true, "payload": resp})
				assert.JSONEq(t, string(respBody), string(data))
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// building stub
			ct := mock.NewMockCategory(ctrl)
			tt.buildStubs(ct)
			// setting stub
			api.App.Category = ct

			recorder := httptest.NewRecorder()
			req, err := http.NewRequest(tt.method, tt.url, tt.body)
			assert.Nil(t, err)
			api.Router.Root.ServeHTTP(recorder, req)
			tt.checkResponse(t, recorder)
		})
	}
}
