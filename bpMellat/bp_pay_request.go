package bpMellat

import (
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type SoapEnvelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    SoapBody `xml:"Body"`
}

type SoapBody struct {
	BpPayRequestResponse BpPayRequestResponse `xml:"bpPayRequestResponse"`
}

type BpPayRequestResponse struct {
	Return string `xml:"return"`
}

func (req *BpMellat) BpPayRequest(input BpPayRequest) (string, string, error) {
	soapEnvelope := `<?xml version="1.0" encoding="UTF-8"?>
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:web="http://interfaces.core.sw.bps.com/">
   <soapenv:Header/>
   <soapenv:Body>
      <web:bpPayRequest>
         <terminalId>` + strconv.Itoa(req.TerminalId) + `</terminalId>
         <userName>` + req.UserName + `</userName>
         <userPassword>` + req.UserPassword + `</userPassword>
         <orderId>` + strconv.FormatInt(input.OrderId, 10) + `</orderId>
         <amount>` + strconv.FormatInt(input.Amount, 10) + `</amount>
         <localDate>` + input.LocalDate + `</localDate>
         <localTime>` + input.LocalTime + `</localTime>
         <additionalData>` + input.AdditionalData + `</additionalData>
         <callBackUrl>` + input.CallBackUrl + `</callBackUrl>
         <payerId>` + strconv.FormatInt(input.PayerId, 10) + `</payerId>
      </web:bpPayRequest>
   </soapenv:Body>
</soapenv:Envelope>`

	resp, err := http.Post(url, "text/xml", bytes.NewBuffer([]byte(soapEnvelope)))
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", err
	}

	var response SoapEnvelope
	err = xml.Unmarshal(body, &response)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(response.Body.BpPayRequestResponse.Return, ",")
	if len(parts) < 1 {
		return "", "", errors.New("پاسخ نامعتبر از سرور بانک")
	}

	if err = getBankError(parts[0]); err != nil {
		return "", "", err
	}

	if len(parts) < 2 {
		return "", "", errors.New("پاسخ نامعتبر از سرور بانک")
	}

	return GatewayURL, parts[1], nil
}
