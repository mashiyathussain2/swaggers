package app

import (
	"encoding/json"
	"fmt"
	"go-app/schema"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/mongo"
)

// Kaleyra contains methods for Kaleyra service functionality
type Kaleyra interface {
	SendOTP(opts *schema.SendOTPOpts) error
}

// KaleyraImpl implements Kaleyra interface methods
type KaleyraImpl struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// KaleyraImplOpts contains args required to create
type KaleyraImplOpts struct {
	App    *App
	DB     *mongo.Database
	Logger *zerolog.Logger
}

// InitKaleyra returns new instance of Kaleyra implementation
func InitKaleyra(opts *KaleyraImplOpts) Kaleyra {
	ui := KaleyraImpl{
		App:    opts.App,
		DB:     opts.DB,
		Logger: opts.Logger,
	}
	return &ui
}

func (ki *KaleyraImpl) SendOTP(opts *schema.SendOTPOpts) error {
	if ki.App.Config.MSGPlatform.Name == "kaleyra" {
		err := ki.SendOTPviaKaleyra(opts)
		if err != nil {
			ki.App.Logger.Err(err).Msg("kaleyra otp error")
			err := ki.SendOTPviaAWS(opts)
			return err
		}
	} else if ki.App.Config.MSGPlatform.Name == "aws" {
		err := ki.SendOTPviaAWS(opts)
		if err != nil {
			ki.App.Logger.Err(err).Msg("AWS otp error")
			err := ki.SendOTPviaKaleyra(opts)
			return err
		}
	}
	return nil
}
func (ki *KaleyraImpl) SendOTPviaAWS(opts *schema.SendOTPOpts) error {
	params := &sns.PublishInput{
		Message:     aws.String(fmt.Sprintf("OTP for login: %s", opts.OTP)),
		PhoneNumber: aws.String(fmt.Sprintf("%s%s", opts.PhoneNo.Prefix, opts.PhoneNo.Number)),
	}
	if _, err := ki.App.SNS.Publish(params); err != nil {
		ki.Logger.Err(err).Interface("phone  no", opts.PhoneNo).Msg("failed to send otp")
		return errors.Wrap(err, "failed to send otp")
	}
	return nil
}
func (ki *KaleyraImpl) SendOTPviaKaleyra(opts *schema.SendOTPOpts) error {

	soURL := "https://api.kaleyra.io/v1/" + ki.App.Config.Kaleyra.SID + "/messages"
	formBody := opts.OTP + " is your OTP for login to Hypd. Watch, Learn & Shop from Lifestyle Videos. #GetHypd"
	data := url.Values{}
	data.Set("to", opts.PhoneNo.Prefix+opts.PhoneNo.Number)
	data.Set("type", "OTP")
	data.Set("sender", "HYPDST")
	data.Set("body", formBody)
	data.Set("template_id", ki.App.Config.Kaleyra.TemplateID)
	client := http.Client{}
	req, err := http.NewRequest(http.MethodPost, soURL, strings.NewReader(data.Encode()))
	if err != nil {
		ki.Logger.Err(err).Interface("request body", data).Msgf("failed to create request to send otp %s", soURL)
		return errors.Wrap(err, "failed to create request to send otp")
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("api-key", ki.App.Config.Kaleyra.APIKey)
	resp, err := client.Do(req)
	//Handle Error
	if err != nil {
		ki.Logger.Err(err).RawJSON("responseBody", []byte(data.Encode())).Msgf("failed to send request to api %s", soURL)
		return errors.Wrap(err, "failed to get response from info")
	}
	defer resp.Body.Close()
	//Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		ki.Logger.Err(err).RawJSON("reqBody", []byte(data.Encode())).Msgf("failed to read response from api %s", soURL)
		return errors.Wrap(err, "failed to get order info")
	}
	var kaleyraResp interface{}
	if err := json.Unmarshal(body, &kaleyraResp); err != nil {
		ki.Logger.Err(err).Str("body", string(body)).Msg("failed to decode body into struct")
		return errors.Wrap(err, "failed to decode body into struct")
	}
	m, ok := (*&kaleyraResp).(map[string]interface{})
	if !ok {
		return fmt.Errorf("want type map[string]interface{};  got %T", opts)
	}
	kError := m["error"]
	if len(kError.(map[string]interface{})) != 0 {
		fmt.Println(kError)
		return errors.New("unable to send otp via kaleyra")
	}
	return nil
}
