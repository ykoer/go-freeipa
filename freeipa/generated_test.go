package freeipa

import (
	"encoding/json"
	"strings"
	"testing"
)

var notFoundIpaResponse = `{
		"result": {
			"result": [],
			"count": 0,
			"truncated": false,
			"summary": "0 services matched"
		},
		"error": null,
		"id": null,
		"principal": "admin@PAAS.REDHAT.COM",
		"version": "4.8.7"
	}`

var actualIpaResponse = `{
		"result": {
			"result": [
				{
					"krbprincipalname": [
						"HTTP/test2.example.com@PAAS.REDHAT.COM"
					],
					"krbcanonicalname": [
						"HTTP/test2.example.com@PAAS.REDHAT.COM"
					],
					"has_keytab": false,
					"dn": "krbprincipalname=HTTP/test2.example.com@PAAS.REDHAT.COM,cn=services, cn=accounts,dc=paas,dc=redhat,dc=com"
				}
			],
			"count": 1,
			"truncated": false,
			"summary": "1 service matched"
		},
		"error": null,
		"id": null,
		"principal": "admin@PAAS.REDHAT.COM",
		"version": "4.8.7"
	}`

var goodIpaResponse = `
	{
		"result": {
			"result": [
				{
					"usercertificate": [
						{
							"__base64__": "MIIE1jCCAz6gAwIBAgIBHzANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9QQUFTLlJFREhBVC5DT00xHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTA1MjExMTQ3NDhaFw0yMTA1MjIxMTQ3NDhaMDYxGDAWBgNVBAoMD1BBQVMuUkVESEFULkNPTTEaMBgGA1UEAwwRdGVzdDIuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDp89peF0ogJo2c0VW6kx0dzV5TtlRKWHdAQRLWK2EONUh/mlC9hpsfK+3IhwlFk90PBcaHDbd8jzAHs8MQF0YF3xrwj+JpughTfO27lJENp+lGvA7S7OHWkdSUyATVHpSBd2dq6qynbp/uzdNIoklfNOfIjs8zWMPX3mk9CMrwH0pL8l1spsV6rDVinJEzucNk8sQAuVKnXI8OJSl5D7hMLjCPBRUSUf5TO2TF5OWZQtHzpOHLR5fRtM55Quaa6w8gUBHzg8AVFGQWd4t18G/cbuwSwdxzglPP7vfOwb7XTMzE+38id/yg8DCGp9/9V0E6MjiuuEYq1yW5j8vu1hLpAgMBAAGjggFpMIIBZTAfBgNVHSMEGDAWgBSXlvS1PLosNp39Gu/uMEJ5vpBOGDBNBggrBgEFBQcBAQRBMD8wPQYIKwYBBQUHMAGGMWh0dHA6Ly9pcGEtY2EuZGV2LmlhZDIuZGMucGFhcy5yZWRoYXQuY29tL2NhL29jc3AwDgYDVR0PAQH/BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBhgYDVR0fBH8wfTB7oEOgQYY/aHR0cDovL2lwYS1jYS5kZXYuaWFkMi5kYy5wYWFzLnJlZGhhdC5jb20vaXBhL2NybC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBTMVPJICCEfprbL0qWy/GyMJFgv1DAcBgNVHREEFTATghF0ZXN0Mi5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAkCO1uUzzo0fgZqgxzUBt7DinG2b7DcHzWGKezMRMhVCVvT79EhHp51SAGoTm3c1v/Yco1KySZVxA5za4QBK/fFkhLlFCULFS4Jab6r0K7fjsKVV48VdJ8hHQM4czzlxKFAb2xTGdaVdXYGjEY6RkEgZF4IO6HL9Uy+Fc36E76+nhDvgVCkEUOaxxZ2Wuy8q5BaP0CzbrUogmdprXcAjYpfTqSzyLKfSrNMbwwpyQ9eJW/slZjuDcays8LPfAO+OXuQGfv6HKwwcPORLMtcQkRaSkpYy9dsB5sG577Vs0SC26E5D0ryZ+40THQRcQic8XrfrQrDRVMMigAyyiLRAIVpIvr9HsHa1cu3sC9hg1+EJ4mC17MZ+00GKu7w3RMzqmrGcCCAyeKnoYzfc9Tuo3lv/my5TX24Cx+sL6hDa5xIC+my6ROscltRw2FHPJP30O7gu3BMnewBFYsUpvInYZP0d5w/57ByoXsqu7XUeMGklZDsmahKhUBW46C8tQxfMW"
						},
						{
							"__base64__": "MIIE1jCCAz6gAwIBAgIBIDANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9QQUFTLlJFREhBVC5DT00xHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTA1MjExMTQ5NDVaFw0yMTA1MjIxMTQ5NDVaMDYxGDAWBgNVBAoMD1BBQVMuUkVESEFULkNPTTEaMBgGA1UEAwwRdGVzdDIuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDh3f9nOfIYUFsSooDa1p8Jf2+m+j9IIqrNAXDeg/PDX95/npPdZJmn3kRaZmDN7SGZG/oElTKG2plSh7T2cyT4eXafqqV1yJvPPiKYZGi0fKJ+X5IwZVd/eUGi60jMDR5ERSWpEG+RFwc+9nl2+9j8rxrhPdd0pLaLsVa0Rfw/KsVy6Zv3rupNDrYEgL8MTywSwX420Ocg8feB449cVk1YHlC6bBIvqlT3M2uB5HYwMlM2XAyylJpb+SCVfKXTscrRd55vSCsS9N1AwLN6R/tEcL7JL4c2VnAueQ0PU8Jt9elS4tzErSXSrkFHmxEzX3+DRbSevJ9P8mZ2t+wblsw/AgMBAAGjggFpMIIBZTAfBgNVHSMEGDAWgBSXlvS1PLosNp39Gu/uMEJ5vpBOGDBNBggrBgEFBQcBAQRBMD8wPQYIKwYBBQUHMAGGMWh0dHA6Ly9pcGEtY2EuZGV2LmlhZDIuZGMucGFhcy5yZWRoYXQuY29tL2NhL29jc3AwDgYDVR0PAQH/BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBhgYDVR0fBH8wfTB7oEOgQYY/aHR0cDovL2lwYS1jYS5kZXYuaWFkMi5kYy5wYWFzLnJlZGhhdC5jb20vaXBhL2NybC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBRq3cgNbsfb4hlIVBk3r57C2o+lBTAcBgNVHREEFTATghF0ZXN0Mi5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAjYkyMyuCptfykUzlZzOu/4SWbRhnTuEzNPEZDU04KeRMkLwQwIYcoCYg7IXp0U++dXrfwdSxSq1L6fmfu1KeHndhU0wbVFz/Rmb9adlWc1joMA9f1Kovcq5lykNwzLb9tQGrAxPX0hvyGI1gjqau0T3Iz1mmIVAfxAH/7Tppk4PqkGe5pBi1n7vYN/aXsVGkT3oWVEbW3jAktFZhx7VN+jFNb5n9wZGOgIAroAfBwizzrKwgr4d6HvArnAyvG0zOv7Wopjmhymp3B+ddLA6CbmwDiFzR2cdkC2quOPu8xAvaPMSYS+pLABaPmLkYSsenn1xNIv51J6DZyiye95JU5mRVCLcENF0xzT4n2Q8p227aWhxDj75PJei8qy3ErcpEB25FR0B1J0rd74i4Er4MXCPorOl6fOAPclmZr7ISOEm9mNCNrx5u37ymGpKgHtPQ1pzUP9d1MHrxL7B54jteuL+pwHgGkwMZNQSc+4ZbBYSjD9AUqlzE4ayICtwSdJ0m"
						}
					],
					"krbcanonicalname": [
						"HTTP/test2.example.com@PAAS.REDHAT.COM"
					],
					"krbprincipalname": [
						"HTTP/test2.example.com@PAAS.REDHAT.COM"
					],
					"has_keytab": false,
					"subject": "CN=test2.example.com,O=PAAS.REDHAT.COM",
					"serial_number": "31",
					"serial_number_hex": "0x1F",
					"issuer": "CN=Certificate Authority,O=PAAS.REDHAT.COM",
					"valid_not_before": "Fri May 21 11:47:48 2021 UTC",
					"valid_not_after": "Sat May 22 11:47:48 2021 UTC",
					"sha1_fingerprint": "bc:01:50:98:ef:c3:b4:28:f0:3a:b8:15:e7:2f:d1:4f:38:0f:53:1e",
					"sha256_fingerprint": "3e:7e:ae:81:97:58:f8:9b:41:aa:bb:4f:5b:d7:43:37:b1:8e:c0:d3:96:1e:5e:25:cc:e6:a0:e9:96:d4:3e:27",
					"dn": "krbprincipalname=HTTP/test2.example.com@PAAS.REDHAT.COM,cn=services,cn=acco* Connection #0 to host ipa.dev.iad2.dc.paas.redhat.com left intactunts,dc=paas,dc=redhat,dc=com"

				}
			],
			"count": 1,
			"truncated": false,
			"summary": "1 service matched"
		},
		"error": null,
		"id": null,
		"principal": "admin@PAAS.REDHAT.COM",
		"version": "4.8.7"
	}`

var goodService = `
	{
		"usercertificate": [
			{
				"__base64__": "MIIE1jCCAz6gAwIBAgIBHzANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9QQUFTLlJFREhBVC5DT00xHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTA1MjExMTQ3NDhaFw0yMTA1MjIxMTQ3NDhaMDYxGDAWBgNVBAoMD1BBQVMuUkVESEFULkNPTTEaMBgGA1UEAwwRdGVzdDIuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDp89peF0ogJo2c0VW6kx0dzV5TtlRKWHdAQRLWK2EONUh/mlC9hpsfK+3IhwlFk90PBcaHDbd8jzAHs8MQF0YF3xrwj+JpughTfO27lJENp+lGvA7S7OHWkdSUyATVHpSBd2dq6qynbp/uzdNIoklfNOfIjs8zWMPX3mk9CMrwH0pL8l1spsV6rDVinJEzucNk8sQAuVKnXI8OJSl5D7hMLjCPBRUSUf5TO2TF5OWZQtHzpOHLR5fRtM55Quaa6w8gUBHzg8AVFGQWd4t18G/cbuwSwdxzglPP7vfOwb7XTMzE+38id/yg8DCGp9/9V0E6MjiuuEYq1yW5j8vu1hLpAgMBAAGjggFpMIIBZTAfBgNVHSMEGDAWgBSXlvS1PLosNp39Gu/uMEJ5vpBOGDBNBggrBgEFBQcBAQRBMD8wPQYIKwYBBQUHMAGGMWh0dHA6Ly9pcGEtY2EuZGV2LmlhZDIuZGMucGFhcy5yZWRoYXQuY29tL2NhL29jc3AwDgYDVR0PAQH/BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBhgYDVR0fBH8wfTB7oEOgQYY/aHR0cDovL2lwYS1jYS5kZXYuaWFkMi5kYy5wYWFzLnJlZGhhdC5jb20vaXBhL2NybC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBTMVPJICCEfprbL0qWy/GyMJFgv1DAcBgNVHREEFTATghF0ZXN0Mi5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAkCO1uUzzo0fgZqgxzUBt7DinG2b7DcHzWGKezMRMhVCVvT79EhHp51SAGoTm3c1v/Yco1KySZVxA5za4QBK/fFkhLlFCULFS4Jab6r0K7fjsKVV48VdJ8hHQM4czzlxKFAb2xTGdaVdXYGjEY6RkEgZF4IO6HL9Uy+Fc36E76+nhDvgVCkEUOaxxZ2Wuy8q5BaP0CzbrUogmdprXcAjYpfTqSzyLKfSrNMbwwpyQ9eJW/slZjuDcays8LPfAO+OXuQGfv6HKwwcPORLMtcQkRaSkpYy9dsB5sG577Vs0SC26E5D0ryZ+40THQRcQic8XrfrQrDRVMMigAyyiLRAIVpIvr9HsHa1cu3sC9hg1+EJ4mC17MZ+00GKu7w3RMzqmrGcCCAyeKnoYzfc9Tuo3lv/my5TX24Cx+sL6hDa5xIC+my6ROscltRw2FHPJP30O7gu3BMnewBFYsUpvInYZP0d5w/57ByoXsqu7XUeMGklZDsmahKhUBW46C8tQxfMW"
			},
			{
				"__base64__": "MIIE1jCCAz6gAwIBAgIBIDANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9QQUFTLlJFREhBVC5DT00xHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTA1MjExMTQ5NDVaFw0yMTA1MjIxMTQ5NDVaMDYxGDAWBgNVBAoMD1BBQVMuUkVESEFULkNPTTEaMBgGA1UEAwwRdGVzdDIuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDh3f9nOfIYUFsSooDa1p8Jf2+m+j9IIqrNAXDeg/PDX95/npPdZJmn3kRaZmDN7SGZG/oElTKG2plSh7T2cyT4eXafqqV1yJvPPiKYZGi0fKJ+X5IwZVd/eUGi60jMDR5ERSWpEG+RFwc+9nl2+9j8rxrhPdd0pLaLsVa0Rfw/KsVy6Zv3rupNDrYEgL8MTywSwX420Ocg8feB449cVk1YHlC6bBIvqlT3M2uB5HYwMlM2XAyylJpb+SCVfKXTscrRd55vSCsS9N1AwLN6R/tEcL7JL4c2VnAueQ0PU8Jt9elS4tzErSXSrkFHmxEzX3+DRbSevJ9P8mZ2t+wblsw/AgMBAAGjggFpMIIBZTAfBgNVHSMEGDAWgBSXlvS1PLosNp39Gu/uMEJ5vpBOGDBNBggrBgEFBQcBAQRBMD8wPQYIKwYBBQUHMAGGMWh0dHA6Ly9pcGEtY2EuZGV2LmlhZDIuZGMucGFhcy5yZWRoYXQuY29tL2NhL29jc3AwDgYDVR0PAQH/BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBhgYDVR0fBH8wfTB7oEOgQYY/aHR0cDovL2lwYS1jYS5kZXYuaWFkMi5kYy5wYWFzLnJlZGhhdC5jb20vaXBhL2NybC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBRq3cgNbsfb4hlIVBk3r57C2o+lBTAcBgNVHREEFTATghF0ZXN0Mi5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAjYkyMyuCptfykUzlZzOu/4SWbRhnTuEzNPEZDU04KeRMkLwQwIYcoCYg7IXp0U++dXrfwdSxSq1L6fmfu1KeHndhU0wbVFz/Rmb9adlWc1joMA9f1Kovcq5lykNwzLb9tQGrAxPX0hvyGI1gjqau0T3Iz1mmIVAfxAH/7Tppk4PqkGe5pBi1n7vYN/aXsVGkT3oWVEbW3jAktFZhx7VN+jFNb5n9wZGOgIAroAfBwizzrKwgr4d6HvArnAyvG0zOv7Wopjmhymp3B+ddLA6CbmwDiFzR2cdkC2quOPu8xAvaPMSYS+pLABaPmLkYSsenn1xNIv51J6DZyiye95JU5mRVCLcENF0xzT4n2Q8p227aWhxDj75PJei8qy3ErcpEB25FR0B1J0rd74i4Er4MXCPorOl6fOAPclmZr7ISOEm9mNCNrx5u37ymGpKgHtPQ1pzUP9d1MHrxL7B54jteuL+pwHgGkwMZNQSc+4ZbBYSjD9AUqlzE4ayICtwSdJ0m"
			}
		],
		"krbcanonicalname": [
			"HTTP/test2.example.com@PAAS.REDHAT.COM"
		],
		"krbprincipalname": [
			"HTTP/test2.example.com@PAAS.REDHAT.COM"
		],
		"has_keytab": false,
		"subject": "CN=test2.example.com,O=PAAS.REDHAT.COM",
		"serial_number": "31",
		"serial_number_hex": "0x1F",
		"issuer": "CN=Certificate Authority,O=PAAS.REDHAT.COM",
		"valid_not_before": "Fri May 21 11:47:48 2021 UTC",
		"valid_not_after": "Sat May 22 11:47:48 2021 UTC",
		"sha1_fingerprint": "bc:01:50:98:ef:c3:b4:28:f0:3a:b8:15:e7:2f:d1:4f:38:0f:53:1e",
		"sha256_fingerprint": "3e:7e:ae:81:97:58:f8:9b:41:aa:bb:4f:5b:d7:43:37:b1:8e:c0:d3:96:1e:5e:25:cc:e6:a0:e9:96:d4:3e:27",
		"dn": "krbprincipalname=HTTP/test2.example.com@PAAS.REDHAT.COM,cn=services,cn=acco* Connection #0 to host ipa.dev.iad2.dc.paas.redhat.com left intactunts,dc=paas,dc=redhat,dc=com"
	}`

var minimalService = `
	{
		"krbcanonicalname": [
			"HTTP/test2.example.com@PAAS.REDHAT.COM"
		],
		"krbprincipalname": [
			"HTTP/test2.example.com@PAAS.REDHAT.COM"
		]
	}`

var badService = `
	{
		"usercertificate": [
			{
				"__base64__": "MIIE1jCCAz6gAwIBAgIBHzANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9QQUFTLlJFREhBVC5DT00xHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTA1MjExMTQ3NDhaFw0yMTA1MjIxMTQ3NDhaMDYxGDAWBgNVBAoMD1BBQVMuUkVESEFULkNPTTEaMBgGA1UEAwwRdGVzdDIuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDp89peF0ogJo2c0VW6kx0dzV5TtlRKWHdAQRLWK2EONUh/mlC9hpsfK+3IhwlFk90PBcaHDbd8jzAHs8MQF0YF3xrwj+JpughTfO27lJENp+lGvA7S7OHWkdSUyATVHpSBd2dq6qynbp/uzdNIoklfNOfIjs8zWMPX3mk9CMrwH0pL8l1spsV6rDVinJEzucNk8sQAuVKnXI8OJSl5D7hMLjCPBRUSUf5TO2TF5OWZQtHzpOHLR5fRtM55Quaa6w8gUBHzg8AVFGQWd4t18G/cbuwSwdxzglPP7vfOwb7XTMzE+38id/yg8DCGp9/9V0E6MjiuuEYq1yW5j8vu1hLpAgMBAAGjggFpMIIBZTAfBgNVHSMEGDAWgBSXlvS1PLosNp39Gu/uMEJ5vpBOGDBNBggrBgEFBQcBAQRBMD8wPQYIKwYBBQUHMAGGMWh0dHA6Ly9pcGEtY2EuZGV2LmlhZDIuZGMucGFhcy5yZWRoYXQuY29tL2NhL29jc3AwDgYDVR0PAQH/BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBhgYDVR0fBH8wfTB7oEOgQYY/aHR0cDovL2lwYS1jYS5kZXYuaWFkMi5kYy5wYWFzLnJlZGhhdC5jb20vaXBhL2NybC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBTMVPJICCEfprbL0qWy/GyMJFgv1DAcBgNVHREEFTATghF0ZXN0Mi5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAkCO1uUzzo0fgZqgxzUBt7DinG2b7DcHzWGKezMRMhVCVvT79EhHp51SAGoTm3c1v/Yco1KySZVxA5za4QBK/fFkhLlFCULFS4Jab6r0K7fjsKVV48VdJ8hHQM4czzlxKFAb2xTGdaVdXYGjEY6RkEgZF4IO6HL9Uy+Fc36E76+nhDvgVCkEUOaxxZ2Wuy8q5BaP0CzbrUogmdprXcAjYpfTqSzyLKfSrNMbwwpyQ9eJW/slZjuDcays8LPfAO+OXuQGfv6HKwwcPORLMtcQkRaSkpYy9dsB5sG577Vs0SC26E5D0ryZ+40THQRcQic8XrfrQrDRVMMigAyyiLRAIVpIvr9HsHa1cu3sC9hg1+EJ4mC17MZ+00GKu7w3RMzqmrGcCCAyeKnoYzfc9Tuo3lv/my5TX24Cx+sL6hDa5xIC+my6ROscltRw2FHPJP30O7gu3BMnewBFYsUpvInYZP0d5w/57ByoXsqu7XUeMGklZDsmahKhUBW46C8tQxfMW"
			},
			{
				"__base64__": "MIIE1jCCAz6gAwIBAgIBIDANBgkqhkiG9w0BAQsFADA6MRgwFgYDVQQKDA9QQUFTLlJFREhBVC5DT00xHjAcBgNVBAMMFUNlcnRpZmljYXRlIEF1dGhvcml0eTAeFw0yMTA1MjExMTQ5NDVaFw0yMTA1MjIxMTQ5NDVaMDYxGDAWBgNVBAoMD1BBQVMuUkVESEFULkNPTTEaMBgGA1UEAwwRdGVzdDIuZXhhbXBsZS5jb20wggEiMA0GCSqGSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDh3f9nOfIYUFsSooDa1p8Jf2+m+j9IIqrNAXDeg/PDX95/npPdZJmn3kRaZmDN7SGZG/oElTKG2plSh7T2cyT4eXafqqV1yJvPPiKYZGi0fKJ+X5IwZVd/eUGi60jMDR5ERSWpEG+RFwc+9nl2+9j8rxrhPdd0pLaLsVa0Rfw/KsVy6Zv3rupNDrYEgL8MTywSwX420Ocg8feB449cVk1YHlC6bBIvqlT3M2uB5HYwMlM2XAyylJpb+SCVfKXTscrRd55vSCsS9N1AwLN6R/tEcL7JL4c2VnAueQ0PU8Jt9elS4tzErSXSrkFHmxEzX3+DRbSevJ9P8mZ2t+wblsw/AgMBAAGjggFpMIIBZTAfBgNVHSMEGDAWgBSXlvS1PLosNp39Gu/uMEJ5vpBOGDBNBggrBgEFBQcBAQRBMD8wPQYIKwYBBQUHMAGGMWh0dHA6Ly9pcGEtY2EuZGV2LmlhZDIuZGMucGFhcy5yZWRoYXQuY29tL2NhL29jc3AwDgYDVR0PAQH/BAQDAgTwMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjCBhgYDVR0fBH8wfTB7oEOgQYY/aHR0cDovL2lwYS1jYS5kZXYuaWFkMi5kYy5wYWFzLnJlZGhhdC5jb20vaXBhL2NybC9NYXN0ZXJDUkwuYmluojSkMjAwMQ4wDAYDVQQKDAVpcGFjYTEeMBwGA1UEAwwVQ2VydGlmaWNhdGUgQXV0aG9yaXR5MB0GA1UdDgQWBBRq3cgNbsfb4hlIVBk3r57C2o+lBTAcBgNVHREEFTATghF0ZXN0Mi5leGFtcGxlLmNvbTANBgkqhkiG9w0BAQsFAAOCAYEAjYkyMyuCptfykUzlZzOu/4SWbRhnTuEzNPEZDU04KeRMkLwQwIYcoCYg7IXp0U++dXrfwdSxSq1L6fmfu1KeHndhU0wbVFz/Rmb9adlWc1joMA9f1Kovcq5lykNwzLb9tQGrAxPX0hvyGI1gjqau0T3Iz1mmIVAfxAH/7Tppk4PqkGe5pBi1n7vYN/aXsVGkT3oWVEbW3jAktFZhx7VN+jFNb5n9wZGOgIAroAfBwizzrKwgr4d6HvArnAyvG0zOv7Wopjmhymp3B+ddLA6CbmwDiFzR2cdkC2quOPu8xAvaPMSYS+pLABaPmLkYSsenn1xNIv51J6DZyiye95JU5mRVCLcENF0xzT4n2Q8p227aWhxDj75PJei8qy3ErcpEB25FR0B1J0rd74i4Er4MXCPorOl6fOAPclmZr7ISOEm9mNCNrx5u37ymGpKgHtPQ1pzUP9d1MHrxL7B54jteuL+pwHgGkwMZNQSc+4ZbBYSjD9AUqlzE4ayICtwSdJ0m"
			}
		],
		"has_keytab": false,
		"subject": "CN=test2.example.com,O=PAAS.REDHAT.COM",
		"serial_number": "31",
		"serial_number_hex": "0x1F",
		"issuer": "CN=Certificate Authority,O=PAAS.REDHAT.COM",
		"valid_not_before": "Fri May 21 11:47:48 2021 UTC",
		"valid_not_after": "Sat May 22 11:47:48 2021 UTC",
		"sha1_fingerprint": "bc:01:50:98:ef:c3:b4:28:f0:3a:b8:15:e7:2f:d1:4f:38:0f:53:1e",
		"sha256_fingerprint": "3e:7e:ae:81:97:58:f8:9b:41:aa:bb:4f:5b:d7:43:37:b1:8e:c0:d3:96:1e:5e:25:cc:e6:a0:e9:96:d4:3e:27",
		"dn": "krbprincipalname=HTTP/test2.example.com@PAAS.REDHAT.COM,cn=services,cn=acco* Connection #0 to host ipa.dev.iad2.dc.paas.redhat.com left intactunts,dc=paas,dc=redhat,dc=com"
	}`

func Test_ServiceFindResultUnmarshal_withNotFoundResponse(t *testing.T) {
	var res serviceFindResponse
	decoder := json.NewDecoder(strings.NewReader(notFoundIpaResponse))
	if e := decoder.Decode(&res); e != nil {
		t.Error(e)
	}
	if res.Result.Count > 0 {
		t.Errorf("We should not have gotten any results.")
	}
}

func Test_ServiceFindResultUnmarshal_withActualResponse(t *testing.T) {
	var res serviceFindResponse
	decoder := json.NewDecoder(strings.NewReader(actualIpaResponse))
	if e := decoder.Decode(&res); e != nil {
		t.Error(e)
	}
}

func Test_ServiceFindResultUnmarshal_withGoodResponse(t *testing.T) {
	var res serviceFindResponse
	decoder := json.NewDecoder(strings.NewReader(goodIpaResponse))
	if e := decoder.Decode(&res); e != nil {
		t.Error(e)
	}
}

func Test_ServicDecode_withGoodResponse(t *testing.T) {
	var res Service
	decoder := json.NewDecoder(strings.NewReader(goodService))
	if e := decoder.Decode(&res); e != nil {
		t.Error(e)
	}
}

func Test_ServiceUnmarshal_withGoodService(t *testing.T) {
	var res Service
	err := json.Unmarshal([]byte(goodService), &res)
	if err != nil {
		t.Error(err)
	}
}

func Test_ServiceUnmarshal_withMinimalService(t *testing.T) {
	var res Service
	err := json.Unmarshal([]byte(minimalService), &res)
	if err != nil {
		t.Error(err)
	}
}

func Test_ServiceUnmarshal_withBadService(t *testing.T) {
	var res Service
	err := json.Unmarshal([]byte(badService), &res)
	if err == nil {
		t.Error(err)
	}
}
