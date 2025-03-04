package config

import (
	"maps"
	"strings"
	"time"

	v1API "github.com/supabase/cli/pkg/api"
	"github.com/supabase/cli/pkg/cast"
	"github.com/supabase/cli/pkg/diff"
)

type (
	auth struct {
		Enabled                bool     `toml:"enabled"`
		Image                  string   `toml:"-"`
		SiteUrl                string   `toml:"site_url"`
		AdditionalRedirectUrls []string `toml:"additional_redirect_urls"`

		JwtExpiry                  uint `toml:"jwt_expiry"`
		EnableRefreshTokenRotation bool `toml:"enable_refresh_token_rotation"`
		RefreshTokenReuseInterval  uint `toml:"refresh_token_reuse_interval"`
		EnableManualLinking        bool `toml:"enable_manual_linking"`

		Hook     hook     `toml:"hook"`
		MFA      mfa      `toml:"mfa"`
		Sessions sessions `toml:"sessions"`

		EnableSignup           bool     `toml:"enable_signup"`
		EnableAnonymousSignIns bool     `toml:"enable_anonymous_sign_ins"`
		Email                  email    `toml:"email"`
		Sms                    sms      `toml:"sms"`
		External               external `toml:"external"`

		// Custom secrets can be injected from .env file
		JwtSecret      string `toml:"-" mapstructure:"jwt_secret"`
		AnonKey        string `toml:"-" mapstructure:"anon_key"`
		ServiceRoleKey string `toml:"-" mapstructure:"service_role_key"`

		ThirdParty thirdParty `toml:"third_party"`
	}

	external map[string]provider

	thirdParty struct {
		Firebase tpaFirebase `toml:"firebase"`
		Auth0    tpaAuth0    `toml:"auth0"`
		Cognito  tpaCognito  `toml:"aws_cognito"`
	}

	tpaFirebase struct {
		Enabled bool `toml:"enabled"`

		ProjectID string `toml:"project_id"`
	}

	tpaAuth0 struct {
		Enabled bool `toml:"enabled"`

		Tenant       string `toml:"tenant"`
		TenantRegion string `toml:"tenant_region"`
	}

	tpaCognito struct {
		Enabled bool `toml:"enabled"`

		UserPoolID     string `toml:"user_pool_id"`
		UserPoolRegion string `toml:"user_pool_region"`
	}

	email struct {
		EnableSignup         bool                     `toml:"enable_signup"`
		DoubleConfirmChanges bool                     `toml:"double_confirm_changes"`
		EnableConfirmations  bool                     `toml:"enable_confirmations"`
		SecurePasswordChange bool                     `toml:"secure_password_change"`
		Template             map[string]emailTemplate `toml:"template"`
		Smtp                 smtp                     `toml:"smtp"`
		MaxFrequency         time.Duration            `toml:"max_frequency"`
		OtpLength            uint                     `toml:"otp_length"`
		OtpExpiry            uint                     `toml:"otp_expiry"`
	}

	smtp struct {
		Host       string `toml:"host"`
		Port       uint16 `toml:"port"`
		User       string `toml:"user"`
		Pass       string `toml:"pass"`
		AdminEmail string `toml:"admin_email"`
		SenderName string `toml:"sender_name"`
	}

	emailTemplate struct {
		Subject     string `toml:"subject"`
		ContentPath string `toml:"content_path"`
	}

	sms struct {
		EnableSignup        bool              `toml:"enable_signup"`
		EnableConfirmations bool              `toml:"enable_confirmations"`
		Template            string            `toml:"template"`
		Twilio              twilioConfig      `toml:"twilio" mapstructure:"twilio"`
		TwilioVerify        twilioConfig      `toml:"twilio_verify" mapstructure:"twilio_verify"`
		Messagebird         messagebirdConfig `toml:"messagebird" mapstructure:"messagebird"`
		Textlocal           textlocalConfig   `toml:"textlocal" mapstructure:"textlocal"`
		Vonage              vonageConfig      `toml:"vonage" mapstructure:"vonage"`
		TestOTP             map[string]string `toml:"test_otp"`
		MaxFrequency        time.Duration     `toml:"max_frequency"`
	}

	hook struct {
		MFAVerificationAttempt      hookConfig `toml:"mfa_verification_attempt"`
		PasswordVerificationAttempt hookConfig `toml:"password_verification_attempt"`
		CustomAccessToken           hookConfig `toml:"custom_access_token"`
		SendSMS                     hookConfig `toml:"send_sms"`
		SendEmail                   hookConfig `toml:"send_email"`
	}

	factorTypeConfiguration struct {
		EnrollEnabled bool `toml:"enroll_enabled"`
		VerifyEnabled bool `toml:"verify_enabled"`
	}

	phoneFactorTypeConfiguration struct {
		factorTypeConfiguration
		OtpLength    uint          `toml:"otp_length"`
		Template     string        `toml:"template"`
		MaxFrequency time.Duration `toml:"max_frequency"`
	}

	mfa struct {
		TOTP               factorTypeConfiguration      `toml:"totp"`
		Phone              phoneFactorTypeConfiguration `toml:"phone"`
		WebAuthn           factorTypeConfiguration      `toml:"web_authn"`
		MaxEnrolledFactors uint                         `toml:"max_enrolled_factors"`
	}

	hookConfig struct {
		Enabled bool   `toml:"enabled"`
		URI     string `toml:"uri"`
		Secrets string `toml:"secrets"`
	}

	sessions struct {
		Timebox           time.Duration `toml:"timebox"`
		InactivityTimeout time.Duration `toml:"inactivity_timeout"`
	}

	twilioConfig struct {
		Enabled           bool   `toml:"enabled"`
		AccountSid        string `toml:"account_sid"`
		MessageServiceSid string `toml:"message_service_sid"`
		AuthToken         string `toml:"auth_token" mapstructure:"auth_token"`
	}

	messagebirdConfig struct {
		Enabled    bool   `toml:"enabled"`
		Originator string `toml:"originator"`
		AccessKey  string `toml:"access_key" mapstructure:"access_key"`
	}

	textlocalConfig struct {
		Enabled bool   `toml:"enabled"`
		Sender  string `toml:"sender"`
		ApiKey  string `toml:"api_key" mapstructure:"api_key"`
	}

	vonageConfig struct {
		Enabled   bool   `toml:"enabled"`
		From      string `toml:"from"`
		ApiKey    string `toml:"api_key" mapstructure:"api_key"`
		ApiSecret string `toml:"api_secret" mapstructure:"api_secret"`
	}

	provider struct {
		Enabled        bool   `toml:"enabled"`
		ClientId       string `toml:"client_id"`
		Secret         string `toml:"secret"`
		Url            string `toml:"url"`
		RedirectUri    string `toml:"redirect_uri"`
		SkipNonceCheck bool   `toml:"skip_nonce_check"`
	}
)

func (a *auth) ToUpdateAuthConfigBody() v1API.UpdateAuthConfigBody {
	body := v1API.UpdateAuthConfigBody{
		SiteUrl:                           &a.SiteUrl,
		UriAllowList:                      cast.Ptr(strings.Join(a.AdditionalRedirectUrls, ",")),
		JwtExp:                            cast.UintToIntPtr(&a.JwtExpiry),
		RefreshTokenRotationEnabled:       &a.EnableRefreshTokenRotation,
		SecurityRefreshTokenReuseInterval: cast.UintToIntPtr(&a.RefreshTokenReuseInterval),
		SecurityManualLinkingEnabled:      &a.EnableManualLinking,
		DisableSignup:                     cast.Ptr(!a.EnableSignup),
		ExternalAnonymousUsersEnabled:     &a.EnableAnonymousSignIns,
	}
	a.Sms.toAuthConfigBody(&body)
	a.External.toAuthConfigBody(&body)
	return body
}

func (a *auth) fromRemoteAuthConfig(remoteConfig v1API.AuthConfigResponse) auth {
	result := *a
	result.SiteUrl = cast.Val(remoteConfig.SiteUrl, "")
	result.AdditionalRedirectUrls = strToArr(cast.Val(remoteConfig.UriAllowList, ""))
	result.JwtExpiry = cast.IntToUint(cast.Val(remoteConfig.JwtExp, 0))
	result.EnableRefreshTokenRotation = cast.Val(remoteConfig.RefreshTokenRotationEnabled, false)
	result.RefreshTokenReuseInterval = cast.IntToUint(cast.Val(remoteConfig.SecurityRefreshTokenReuseInterval, 0))
	result.EnableManualLinking = cast.Val(remoteConfig.SecurityManualLinkingEnabled, false)
	result.EnableSignup = !cast.Val(remoteConfig.DisableSignup, false)
	result.EnableAnonymousSignIns = cast.Val(remoteConfig.ExternalAnonymousUsersEnabled, false)
	result.Sms.fromAuthConfig(remoteConfig)
	result.External = maps.Clone(result.External)
	result.External.fromAuthConfig(remoteConfig)
	return result
}

func (s sms) toAuthConfigBody(body *v1API.UpdateAuthConfigBody) {
	body.ExternalPhoneEnabled = &s.EnableSignup
	body.SmsMaxFrequency = cast.Ptr(int(s.MaxFrequency.Seconds()))
	body.SmsAutoconfirm = &s.EnableConfirmations
	body.SmsTemplate = &s.Template
	if otpString := mapToEnv(s.TestOTP); len(otpString) > 0 {
		body.SmsTestOtp = &otpString
		// Set a 10 year validity for test OTP
		timestamp := time.Now().UTC().AddDate(10, 0, 0).Format(time.RFC3339)
		body.SmsTestOtpValidUntil = &timestamp
	}
	// Api only overrides configs of enabled providers
	switch {
	case s.Twilio.Enabled:
		body.SmsProvider = cast.Ptr("twilio")
		body.SmsTwilioAuthToken = &s.Twilio.AuthToken
		body.SmsTwilioAccountSid = &s.Twilio.AccountSid
		body.SmsTwilioMessageServiceSid = &s.Twilio.MessageServiceSid
	case s.TwilioVerify.Enabled:
		body.SmsProvider = cast.Ptr("twilio_verify")
		body.SmsTwilioVerifyAuthToken = &s.TwilioVerify.AuthToken
		body.SmsTwilioVerifyAccountSid = &s.TwilioVerify.AccountSid
		body.SmsTwilioVerifyMessageServiceSid = &s.TwilioVerify.MessageServiceSid
	case s.Messagebird.Enabled:
		body.SmsProvider = cast.Ptr("messagebird")
		body.SmsMessagebirdAccessKey = &s.Messagebird.AccessKey
		body.SmsMessagebirdOriginator = &s.Messagebird.Originator
	case s.Textlocal.Enabled:
		body.SmsProvider = cast.Ptr("textlocal")
		body.SmsTextlocalApiKey = &s.Textlocal.ApiKey
		body.SmsTextlocalSender = &s.Textlocal.Sender
	case s.Vonage.Enabled:
		body.SmsProvider = cast.Ptr("vonage")
		body.SmsVonageApiSecret = &s.Vonage.ApiSecret
		body.SmsVonageApiKey = &s.Vonage.ApiKey
		body.SmsVonageFrom = &s.Vonage.From
	}
}

func (s *sms) fromAuthConfig(remoteConfig v1API.AuthConfigResponse) {
	s.EnableSignup = cast.Val(remoteConfig.ExternalPhoneEnabled, false)
	s.MaxFrequency = time.Duration(cast.Val(remoteConfig.SmsMaxFrequency, 0)) * time.Second
	s.EnableConfirmations = cast.Val(remoteConfig.SmsAutoconfirm, false)
	s.Template = cast.Val(remoteConfig.SmsTemplate, "")
	s.TestOTP = envToMap(cast.Val(remoteConfig.SmsTestOtp, ""))
	// We are only interested in the provider that's enabled locally
	switch {
	case s.Twilio.Enabled:
		s.Twilio.AuthToken = hashPrefix + cast.Val(remoteConfig.SmsTwilioAuthToken, "")
		s.Twilio.AccountSid = cast.Val(remoteConfig.SmsTwilioAccountSid, "")
		s.Twilio.MessageServiceSid = cast.Val(remoteConfig.SmsTwilioMessageServiceSid, "")
	case s.TwilioVerify.Enabled:
		s.TwilioVerify.AuthToken = hashPrefix + cast.Val(remoteConfig.SmsTwilioVerifyAuthToken, "")
		s.TwilioVerify.AccountSid = cast.Val(remoteConfig.SmsTwilioVerifyAccountSid, "")
		s.TwilioVerify.MessageServiceSid = cast.Val(remoteConfig.SmsTwilioVerifyMessageServiceSid, "")
	case s.Messagebird.Enabled:
		s.Messagebird.AccessKey = hashPrefix + cast.Val(remoteConfig.SmsMessagebirdAccessKey, "")
		s.Messagebird.Originator = cast.Val(remoteConfig.SmsMessagebirdOriginator, "")
	case s.Textlocal.Enabled:
		s.Textlocal.ApiKey = hashPrefix + cast.Val(remoteConfig.SmsTextlocalApiKey, "")
		s.Textlocal.Sender = cast.Val(remoteConfig.SmsTextlocalSender, "")
	case s.Vonage.Enabled:
		s.Vonage.ApiSecret = hashPrefix + cast.Val(remoteConfig.SmsVonageApiSecret, "")
		s.Vonage.ApiKey = cast.Val(remoteConfig.SmsVonageApiKey, "")
		s.Vonage.From = cast.Val(remoteConfig.SmsVonageFrom, "")
	case !s.EnableSignup:
		// Nothing to do if both local and remote providers are disabled.
		return
	}
	if provider := cast.Val(remoteConfig.SmsProvider, ""); len(provider) > 0 {
		s.Twilio.Enabled = provider == "twilio"
		s.TwilioVerify.Enabled = provider == "twilio_verify"
		s.Messagebird.Enabled = provider == "messagebird"
		s.Textlocal.Enabled = provider == "textlocal"
		s.Vonage.Enabled = provider == "vonage"
	}
}

func (e external) toAuthConfigBody(body *v1API.UpdateAuthConfigBody) {
	if len(e) == 0 {
		return
	}
	var p *provider
	// Ignore configs of disabled providers because their envs are not loaded
	p = cast.Ptr(e["apple"])
	if body.ExternalAppleEnabled = &p.Enabled; *body.ExternalAppleEnabled {
		body.ExternalAppleClientId = &p.ClientId
		body.ExternalAppleSecret = &p.Secret
	}
	p = cast.Ptr(e["azure"])
	if body.ExternalAzureEnabled = &p.Enabled; *body.ExternalAzureEnabled {
		body.ExternalAzureClientId = &p.ClientId
		body.ExternalAzureSecret = &p.Secret
		body.ExternalAzureUrl = &p.Url
	}
	p = cast.Ptr(e["bitbucket"])
	if body.ExternalBitbucketEnabled = &p.Enabled; *body.ExternalBitbucketEnabled {
		body.ExternalBitbucketClientId = &p.ClientId
		body.ExternalBitbucketSecret = &p.Secret
	}
	p = cast.Ptr(e["discord"])
	if body.ExternalDiscordEnabled = &p.Enabled; *body.ExternalDiscordEnabled {
		body.ExternalDiscordClientId = &p.ClientId
		body.ExternalDiscordSecret = &p.Secret
	}
	p = cast.Ptr(e["facebook"])
	if body.ExternalFacebookEnabled = &p.Enabled; *body.ExternalFacebookEnabled {
		body.ExternalFacebookClientId = &p.ClientId
		body.ExternalFacebookSecret = &p.Secret
	}
	p = cast.Ptr(e["figma"])
	if body.ExternalFigmaEnabled = &p.Enabled; *body.ExternalFigmaEnabled {
		body.ExternalFigmaClientId = &p.ClientId
		body.ExternalFigmaSecret = &p.Secret
	}
	p = cast.Ptr(e["github"])
	if body.ExternalGithubEnabled = &p.Enabled; *body.ExternalGithubEnabled {
		body.ExternalGithubClientId = &p.ClientId
		body.ExternalGithubSecret = &p.Secret
	}
	p = cast.Ptr(e["gitlab"])
	if body.ExternalGitlabEnabled = &p.Enabled; *body.ExternalGitlabEnabled {
		body.ExternalGitlabClientId = &p.ClientId
		body.ExternalGitlabSecret = &p.Secret
		body.ExternalGitlabUrl = &p.Url
	}
	p = cast.Ptr(e["google"])
	if body.ExternalGoogleEnabled = &p.Enabled; *body.ExternalGoogleEnabled {
		body.ExternalGoogleClientId = &p.ClientId
		body.ExternalGoogleSecret = &p.Secret
		body.ExternalGoogleSkipNonceCheck = &p.SkipNonceCheck
	}
	p = cast.Ptr(e["kakao"])
	if body.ExternalKakaoEnabled = &p.Enabled; *body.ExternalKakaoEnabled {
		body.ExternalKakaoClientId = &p.ClientId
		body.ExternalKakaoSecret = &p.Secret
	}
	p = cast.Ptr(e["keycloak"])
	if body.ExternalKeycloakEnabled = &p.Enabled; *body.ExternalKeycloakEnabled {
		body.ExternalKeycloakClientId = &p.ClientId
		body.ExternalKeycloakSecret = &p.Secret
		body.ExternalKeycloakUrl = &p.Url
	}
	p = cast.Ptr(e["linkedin_oidc"])
	if body.ExternalLinkedinOidcEnabled = &p.Enabled; *body.ExternalLinkedinOidcEnabled {
		body.ExternalLinkedinOidcClientId = &p.ClientId
		body.ExternalLinkedinOidcSecret = &p.Secret
	}
	p = cast.Ptr(e["notion"])
	if body.ExternalNotionEnabled = &p.Enabled; *body.ExternalNotionEnabled {
		body.ExternalNotionClientId = &p.ClientId
		body.ExternalNotionSecret = &p.Secret
	}
	p = cast.Ptr(e["slack_oidc"])
	if body.ExternalSlackOidcEnabled = &p.Enabled; *body.ExternalSlackOidcEnabled {
		body.ExternalSlackOidcClientId = &p.ClientId
		body.ExternalSlackOidcSecret = &p.Secret
	}
	p = cast.Ptr(e["spotify"])
	if body.ExternalSpotifyEnabled = &p.Enabled; *body.ExternalSpotifyEnabled {
		body.ExternalSpotifyClientId = &p.ClientId
		body.ExternalSpotifySecret = &p.Secret
	}
	p = cast.Ptr(e["twitch"])
	if body.ExternalTwitchEnabled = &p.Enabled; *body.ExternalTwitchEnabled {
		body.ExternalTwitchClientId = &p.ClientId
		body.ExternalTwitchSecret = &p.Secret
	}
	p = cast.Ptr(e["twitter"])
	if body.ExternalTwitterEnabled = &p.Enabled; *body.ExternalTwitterEnabled {
		body.ExternalTwitterClientId = &p.ClientId
		body.ExternalTwitterSecret = &p.Secret
	}
	p = cast.Ptr(e["workos"])
	if body.ExternalWorkosEnabled = &p.Enabled; *body.ExternalWorkosEnabled {
		body.ExternalWorkosClientId = &p.ClientId
		body.ExternalWorkosSecret = &p.Secret
		body.ExternalWorkosUrl = &p.Url
	}
	p = cast.Ptr(e["zoom"])
	if body.ExternalZoomEnabled = &p.Enabled; *body.ExternalZoomEnabled {
		body.ExternalZoomClientId = &p.ClientId
		body.ExternalZoomSecret = &p.Secret
	}
}

func (e external) fromAuthConfig(remoteConfig v1API.AuthConfigResponse) {
	if len(e) == 0 {
		return
	}
	var p provider
	// Ignore configs of disabled providers because their envs are not loaded
	if p = e["apple"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalAppleClientId, "")
		if ids := cast.Val(remoteConfig.ExternalAppleAdditionalClientIds, ""); len(ids) > 0 {
			p.ClientId += "," + ids
		}
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalAppleSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalAppleEnabled, false)
	e["apple"] = p

	if p = e["azure"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalAzureClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalAzureSecret, "")
		p.Url = cast.Val(remoteConfig.ExternalAzureUrl, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalAzureEnabled, false)
	e["azure"] = p

	if p = e["bitbucket"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalBitbucketClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalBitbucketSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalBitbucketEnabled, false)
	e["bitbucket"] = p

	if p = e["discord"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalDiscordClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalDiscordSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalDiscordEnabled, false)
	e["discord"] = p

	if p = e["facebook"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalFacebookClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalFacebookSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalFacebookEnabled, false)
	e["facebook"] = p

	if p = e["figma"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalFigmaClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalFigmaSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalFigmaEnabled, false)
	e["figma"] = p

	if p = e["github"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalGithubClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalGithubSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalGithubEnabled, false)
	e["github"] = p

	if p = e["gitlab"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalGitlabClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalGitlabSecret, "")
		p.Url = cast.Val(remoteConfig.ExternalGitlabUrl, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalGitlabEnabled, false)
	e["gitlab"] = p

	if p = e["google"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalGoogleClientId, "")
		if ids := cast.Val(remoteConfig.ExternalGoogleAdditionalClientIds, ""); len(ids) > 0 {
			p.ClientId += "," + ids
		}
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalGoogleSecret, "")
		p.SkipNonceCheck = cast.Val(remoteConfig.ExternalGoogleSkipNonceCheck, false)
	}
	p.Enabled = cast.Val(remoteConfig.ExternalGoogleEnabled, false)
	e["google"] = p

	if p = e["kakao"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalKakaoClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalKakaoSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalKakaoEnabled, false)
	e["kakao"] = p

	if p = e["keycloak"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalKeycloakClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalKeycloakSecret, "")
		p.Url = cast.Val(remoteConfig.ExternalKeycloakUrl, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalKeycloakEnabled, false)
	e["keycloak"] = p

	if p = e["linkedin_oidc"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalLinkedinOidcClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalLinkedinOidcSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalLinkedinOidcEnabled, false)
	e["linkedin_oidc"] = p

	if p = e["notion"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalNotionClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalNotionSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalNotionEnabled, false)
	e["notion"] = p

	if p = e["slack_oidc"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalSlackOidcClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalSlackOidcSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalSlackOidcEnabled, false)
	e["slack_oidc"] = p

	if p = e["spotify"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalSpotifyClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalSpotifySecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalSpotifyEnabled, false)
	e["spotify"] = p

	if p = e["twitch"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalTwitchClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalTwitchSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalTwitchEnabled, false)
	e["twitch"] = p

	if p = e["twitter"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalTwitterClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalTwitterSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalTwitterEnabled, false)
	e["twitter"] = p

	if p = e["workos"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalWorkosClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalWorkosSecret, "")
		p.Url = cast.Val(remoteConfig.ExternalWorkosUrl, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalWorkosEnabled, false)
	e["workos"] = p

	if p = e["zoom"]; p.Enabled {
		p.ClientId = cast.Val(remoteConfig.ExternalZoomClientId, "")
		p.Secret = hashPrefix + cast.Val(remoteConfig.ExternalZoomSecret, "")
	}
	p.Enabled = cast.Val(remoteConfig.ExternalZoomEnabled, false)
	e["zoom"] = p
}

func (a *auth) DiffWithRemote(projectRef string, remoteConfig v1API.AuthConfigResponse) ([]byte, error) {
	hashed := a.hashSecrets(projectRef)
	// Convert the config values into easily comparable remoteConfig values
	currentValue, err := ToTomlBytes(hashed)
	if err != nil {
		return nil, err
	}
	remoteCompare, err := ToTomlBytes(hashed.fromRemoteAuthConfig(remoteConfig))
	if err != nil {
		return nil, err
	}
	return diff.Diff("remote[auth]", remoteCompare, "local[auth]", currentValue), nil
}

const hashPrefix = "hash:"

func (a *auth) hashSecrets(key string) auth {
	hash := func(v string) string {
		return hashPrefix + sha256Hmac(key, v)
	}
	result := *a
	if len(result.Email.Smtp.Pass) > 0 {
		result.Email.Smtp.Pass = hash(result.Email.Smtp.Pass)
	}
	// Only hash secrets for locally enabled providers because other envs won't be loaded
	switch {
	case result.Sms.Twilio.Enabled:
		result.Sms.Twilio.AuthToken = hash(result.Sms.Twilio.AuthToken)
	case result.Sms.TwilioVerify.Enabled:
		result.Sms.TwilioVerify.AuthToken = hash(result.Sms.TwilioVerify.AuthToken)
	case result.Sms.Messagebird.Enabled:
		result.Sms.Messagebird.AccessKey = hash(result.Sms.Messagebird.AccessKey)
	case result.Sms.Textlocal.Enabled:
		result.Sms.Textlocal.ApiKey = hash(result.Sms.Textlocal.ApiKey)
	case result.Sms.Vonage.Enabled:
		result.Sms.Vonage.ApiSecret = hash(result.Sms.Vonage.ApiSecret)
	}
	if result.Hook.MFAVerificationAttempt.Enabled {
		result.Hook.MFAVerificationAttempt.Secrets = hash(result.Hook.MFAVerificationAttempt.Secrets)
	}
	if result.Hook.PasswordVerificationAttempt.Enabled {
		result.Hook.PasswordVerificationAttempt.Secrets = hash(result.Hook.PasswordVerificationAttempt.Secrets)
	}
	if result.Hook.CustomAccessToken.Enabled {
		result.Hook.CustomAccessToken.Secrets = hash(result.Hook.CustomAccessToken.Secrets)
	}
	if result.Hook.SendSMS.Enabled {
		result.Hook.SendSMS.Secrets = hash(result.Hook.SendSMS.Secrets)
	}
	if result.Hook.SendEmail.Enabled {
		result.Hook.SendEmail.Secrets = hash(result.Hook.SendEmail.Secrets)
	}
	if size := len(a.External); size > 0 {
		result.External = make(map[string]provider, size)
	}
	for name, provider := range a.External {
		if provider.Enabled {
			provider.Secret = hash(provider.Secret)
		}
		result.External[name] = provider
	}
	// TODO: support SecurityCaptchaSecret in local config
	return result
}
