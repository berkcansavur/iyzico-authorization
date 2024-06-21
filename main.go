package iyzico

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

func GenerateAuthorizationAndPkiString(apiKey string, secretKey string, request InitializeBkmRequest, rnd string) (string, string, error) {
	if apiKey == "" || secretKey == "" || rnd == "" {
		return "", "", errors.New("apiKey, secretKey, or rnd is empty")
	}
	if err := ValidateInitializeBkmRequest(request); err != nil {
		return "", "", err
	}

	requestString := FormatInitializeBkm(request)
	return GetAuthorizationAndPkiString(apiKey, rnd, secretKey, requestString)
}

func GenerateAuthorizationAndPkiStringForCreatePayment(apiKey string, secretKey string, request CreatePaymentRequest, rnd string) (string, string, error) {
	if apiKey == "" || secretKey == "" || rnd == "" {
		return "", "", errors.New("apiKey, secretKey, or rnd is empty")
	}
	if err := ValidateCreatePaymentRequest(request); err != nil {
		return "", "", err
	}

	requestString := FormatCreatePayment(request)
	return GetAuthorizationAndPkiString(apiKey, rnd, secretKey, requestString)
}

func GetAuthorizationAndPkiString(apiKey string, rnd string, secretKey string, requestString string) (string, string, error) {
	if apiKey == "" || rnd == "" || secretKey == "" || requestString == "" {
		return "", "", errors.New("apiKey, rnd, secretKey, or requestString is empty")
	}

	hash := sha1.New()
	hash.Write([]byte(apiKey + rnd + secretKey + requestString))
	hashInBase64 := base64.StdEncoding.EncodeToString(hash.Sum(nil))
	authorization := fmt.Sprintf("IYZWS %s:%s", apiKey, hashInBase64)
	pkiString := apiKey + rnd + secretKey + requestString
	return authorization, pkiString, nil
}

type CreatePaymentRequest struct {
	Locale          string         `json:"locale"`
	ConversationID  string         `json:"conversationId"`
	Price           string         `json:"price"`
	PaidPrice       string         `json:"paidPrice"`
	Installment     int            `json:"installment"`
	PaymentChannel  string         `json:"paymentChannel"`
	BasketID        string         `json:"basketId"`
	PaymentGroup    string         `json:"paymentGroup"`
	PaymentCard     PaymentCard    `json:"paymentCard"`
	Buyer           Buyer          `json:"buyer"`
	ShippingAddress Address        `json:"shippingAddress"`
	BillingAddress  BillingAddress `json:"billingAddress"`
	BasketItems     []BasketItem   `json:"basketItems"`
	Currency        string         `json:"currency"`
}

type InitializeBkmRequest struct {
	Locale          string       `json:"locale,omitempty"`
	ConversationID  string       `json:"conversationId,omitempty"`
	Price           string       `json:"price,omitempty"`
	PaymentChannel  string       `json:"paymentChannel,omitempty"`
	BasketID        string       `json:"basketId,omitempty"`
	PaymentGroup    string       `json:"paymentGroup,omitempty"`
	PaymentCard     PaymentCard  `json:"paymentCard,omitempty"`
	Buyer           Buyer        `json:"buyer,omitempty"`
	ShippingAddress Address      `json:"shippingAddress,omitempty"`
	BillingAddress  Address      `json:"billingAddress,omitempty"`
	BasketItems     []BasketItem `json:"basketItems,omitempty"`
	CallbackURL     string       `json:"callbackUrl,omitempty"`
}
type PaymentCard struct {
	CardHolderName string `json:"cardHolderName"`
	CardNumber     string `json:"cardNumber"`
	ExpireYear     string `json:"expireYear"`
	ExpireMonth    string `json:"expireMonth"`
	Cvc            string `json:"cvc"`
	RegisterCard   int    `json:"registerCard"`
}
type Buyer struct {
	Id                  string `json:"id"`
	Name                string `json:"name"`
	Surname             string `json:"surname"`
	IdentityNumber      string `json:"identityNumber"`
	Email               string `json:"email"`
	GsmNumber           string `json:"gsmNumber"`
	RegistrationDate    string `json:"registrationDate"`
	LastLoginDate       string `json:"lastLoginDate"`
	RegistrationAddress string `json:"registrationAddress"`
	City                string `json:"city"`
	Country             string `json:"country"`
	ZipCode             string `json:"zipCode"`
	Ip                  string `json:"ip"`
}
type Address struct {
	Address     string `json:"address"`
	ZipCode     string `json:"zipCode"`
	ContactName string `json:"contactName"`
	City        string `json:"city"`
	Country     string `json:"country"`
}
type BasketItem struct {
	Id        string `json:"id"`
	Name      string `json:"name"`
	Category1 string `json:"category1"`
	Category2 string `json:"category2"`
	ItemType  string `json:"itemType"`
	Price     string `json:"price"`
}

type BillingAddress struct {
	Address     string `json:"address"`
	ContactName string `json:"contactName"`
	City        string `json:"city"`
	Country     string `json:"country"`
}

func FormatBasketItems(items []BasketItem) string {
	var result []string
	for _, item := range items {
		itemStr := fmt.Sprintf("[id=%s,price=%s,name=%s,category1=%s,category2=%s,itemType=%s]",
			item.Id, item.Price, item.Name, item.Category1, item.Category2, item.ItemType)
		result = append(result, itemStr)
	}
	return strings.Join(result, ", ")
}
func FormatAddress(address Address) string {
	return fmt.Sprintf("[address=%s,zipCode=%s,contactName=%s,city=%s,country=%s]",
		address.Address, address.ZipCode, address.ContactName, address.City, address.Country)
}
func FormatBuyer(buyer Buyer) string {
	return fmt.Sprintf("[id=%s,name=%s,surname=%s,identityNumber=%s,email=%s,gsmNumber=%s,registrationDate=%s,lastLoginDate=%s,registrationAddress=%s,city=%s,country=%s,zipCode=%s,ip=%s]",
		buyer.Id, buyer.Name, buyer.Surname, buyer.IdentityNumber, buyer.Email, buyer.GsmNumber, buyer.RegistrationDate, buyer.LastLoginDate, buyer.RegistrationAddress, buyer.City, buyer.Country, buyer.ZipCode, buyer.Ip)
}

func FormatInitializeBkm(initializeBkm InitializeBkmRequest) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(fmt.Sprintf("locale=%s,", initializeBkm.Locale))
	buffer.WriteString(fmt.Sprintf("conversationId=%s,", initializeBkm.ConversationID))
	buffer.WriteString(fmt.Sprintf("price=%s,", initializeBkm.Price))
	buffer.WriteString(fmt.Sprintf("basketId=%s,", initializeBkm.BasketID))
	buffer.WriteString(fmt.Sprintf("paymentGroup=%s,", initializeBkm.PaymentGroup))
	buffer.WriteString("buyer=")
	buffer.WriteString(FormatBuyer(initializeBkm.Buyer))
	buffer.WriteString(",")
	buffer.WriteString("shippingAddress=")
	buffer.WriteString(FormatAddress(initializeBkm.ShippingAddress))
	buffer.WriteString(",")
	buffer.WriteString("billingAddress=")
	buffer.WriteString(FormatAddress(initializeBkm.BillingAddress))
	buffer.WriteString(",")
	buffer.WriteString("basketItems=[")
	buffer.WriteString(FormatBasketItems(initializeBkm.BasketItems))
	buffer.WriteString("],")
	buffer.WriteString(fmt.Sprintf("callbackUrl=%s", initializeBkm.CallbackURL))
	buffer.WriteString("]")

	return buffer.String()
}
func FormatPaymentCard(paymentCard PaymentCard) string {
	return fmt.Sprintf("[cardHolderName=%s,cardNumber=%s,expireYear=%s,expireMonth=%s,cvc=%s,registerCard=%d]",
		paymentCard.CardHolderName, paymentCard.CardNumber, paymentCard.ExpireYear, paymentCard.ExpireMonth, paymentCard.Cvc, paymentCard.RegisterCard)
}
func FormatBillingAddress(address BillingAddress) string {
	return fmt.Sprintf("[address=%s,contactName=%s,city=%s,country=%s]",
		address.Address, address.ContactName, address.City, address.Country)
}
func FormatCreatePayment(createPayment CreatePaymentRequest) string {
	var buffer bytes.Buffer
	buffer.WriteString("[")
	buffer.WriteString(fmt.Sprintf("locale=%s,", createPayment.Locale))
	buffer.WriteString(fmt.Sprintf("conversationId=%s,", createPayment.ConversationID))
	buffer.WriteString(fmt.Sprintf("price=%s,", createPayment.Price))
	buffer.WriteString(fmt.Sprintf("paidPrice=%s,", createPayment.PaidPrice))
	buffer.WriteString(fmt.Sprintf("installment=%d,", createPayment.Installment))
	buffer.WriteString(fmt.Sprintf("paymentChannel=%s,", createPayment.PaymentChannel))
	buffer.WriteString(fmt.Sprintf("basketId=%s,", createPayment.BasketID))
	buffer.WriteString(fmt.Sprintf("paymentGroup=%s,", createPayment.PaymentGroup))
	buffer.WriteString("paymentCard=")
	buffer.WriteString(FormatPaymentCard(createPayment.PaymentCard))
	buffer.WriteString(",")
	buffer.WriteString("buyer=")
	buffer.WriteString(FormatBuyer(createPayment.Buyer))
	buffer.WriteString(",")
	buffer.WriteString("shippingAddress=")
	buffer.WriteString(FormatAddress(createPayment.ShippingAddress))
	buffer.WriteString(",")
	buffer.WriteString("billingAddress=")
	buffer.WriteString(FormatBillingAddress(createPayment.BillingAddress))
	buffer.WriteString(",")
	buffer.WriteString("basketItems=[")
	buffer.WriteString(FormatBasketItems(createPayment.BasketItems))
	buffer.WriteString("],")
	buffer.WriteString(fmt.Sprintf("currency=%s", createPayment.Currency))
	buffer.WriteString("]")

	return buffer.String()
}
func ValidateInitializeBkmRequest(request InitializeBkmRequest) error {
	if request.Locale == "" {
		return errors.New("Locale is empty")
	}
	if request.ConversationID == "" {
		return errors.New("ConversationID is empty")
	}
	if request.Price == "" {
		return errors.New("Price is empty")
	}
	if request.PaymentChannel == "" {
		return errors.New("PaymentChannel is empty")
	}
	if request.BasketID == "" {
		return errors.New("BasketID is empty")
	}
	if request.PaymentGroup == "" {
		return errors.New("PaymentGroup is empty")
	}
	if err := ValidatePaymentCard(request.PaymentCard); err != nil {
		return err
	}
	if err := ValidateBuyer(request.Buyer); err != nil {
		return err
	}
	if err := ValidateAddress(request.ShippingAddress); err != nil {
		return err
	}
	if err := ValidateAddress(request.BillingAddress); err != nil {
		return err
	}
	if len(request.BasketItems) == 0 {
		return errors.New("BasketItems is empty")
	}
	if request.CallbackURL == "" {
		return errors.New("CallbackURL is empty")
	}
	return nil
}
func ValidateCreatePaymentRequest(request CreatePaymentRequest) error {
	if request.Locale == "" {
		return errors.New("Locale is empty")
	}
	if request.ConversationID == "" {
		return errors.New("ConversationID is empty")
	}
	if request.Price == "" {
		return errors.New("Price is empty")
	}
	if request.PaidPrice == "" {
		return errors.New("PaidPrice is empty")
	}
	if request.Installment == 0 {
		return errors.New("Installment is zero")
	}
	if request.PaymentChannel == "" {
		return errors.New("PaymentChannel is empty")
	}
	if request.BasketID == "" {
		return errors.New("BasketID is empty")
	}
	if request.PaymentGroup == "" {
		return errors.New("PaymentGroup is empty")
	}
	if err := ValidatePaymentCard(request.PaymentCard); err != nil {
		return err
	}
	if err := ValidateBuyer(request.Buyer); err != nil {
		return err
	}
	if err := ValidateAddress(request.ShippingAddress); err != nil {
		return err
	}
	if err := ValidateBillingAddress(request.BillingAddress); err != nil {
		return err
	}
	if len(request.BasketItems) == 0 {
		return errors.New("BasketItems is empty")
	}
	if request.Currency == "" {
		return errors.New("Currency is empty")
	}
	return nil
}

func ValidatePaymentCard(card PaymentCard) error {
	if card.CardHolderName == "" {
		return errors.New("CardHolderName is empty")
	}
	if card.CardNumber == "" {
		return errors.New("CardNumber is empty")
	}
	if card.ExpireYear == "" {
		return errors.New("ExpireYear is empty")
	}
	if card.ExpireMonth == "" {
		return errors.New("ExpireMonth is empty")
	}
	if card.Cvc == "" {
		return errors.New("Cvc is empty")
	}
	if card.RegisterCard == 0 {
		return errors.New("RegisterCard is zero")
	}
	return nil
}

func ValidateBuyer(buyer Buyer) error {
	if buyer.Id == "" {
		return errors.New("Id is empty")
	}
	if buyer.Name == "" {
		return errors.New("Name is empty")
	}
	if buyer.Surname == "" {
		return errors.New("Surname is empty")
	}
	if buyer.IdentityNumber == "" {
		return errors.New("IdentityNumber is empty")
	}
	if buyer.Email == "" {
		return errors.New("Email is empty")
	}
	if buyer.GsmNumber == "" {
		return errors.New("GsmNumber is empty")
	}
	if buyer.RegistrationDate == "" {
		return errors.New("RegistrationDate is empty")
	}
	if buyer.LastLoginDate == "" {
		return errors.New("LastLoginDate is empty")
	}
	if buyer.RegistrationAddress == "" {
		return errors.New("RegistrationAddress is empty")
	}
	if buyer.City == "" {
		return errors.New("City is empty")
	}
	if buyer.Country == "" {
		return errors.New("Country is empty")
	}
	if buyer.ZipCode == "" {
		return errors.New("ZipCode is empty")
	}
	if buyer.Ip == "" {
		return errors.New("Ip is empty")
	}
	return nil
}

func ValidateAddress(address Address) error {
	if address.Address == "" {
		return errors.New("Address is empty")
	}
	if address.ZipCode == "" {
		return errors.New("ZipCode is empty")
	}
	if address.ContactName == "" {
		return errors.New("ContactName is empty")
	}
	if address.City == "" {
		return errors.New("City is empty")
	}
	if address.Country == "" {
		return errors.New("Country is empty")
	}
	return nil
}

func ValidateBillingAddress(address BillingAddress) error {
	if address.Address == "" {
		return errors.New("Address is empty")
	}
	if address.ContactName == "" {
		return errors.New("ContactName is empty")
	}
	if address.City == "" {
		return errors.New("City is empty")
	}
	if address.Country == "" {
		return errors.New("Country is empty")
	}
	return nil
}
