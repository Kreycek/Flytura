// config/config.go
package flytura

import (
	"Flytura/internal/models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"

	"cloud.google.com/go/bigquery"
	"cloud.google.com/go/civil"
	"github.com/extrame/xls"
	"github.com/golang-jwt/jwt"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/text/unicode/norm"
)

// Variável global que contém a chave secreta para JWT

var SecretKey = []byte("my_secret_key")
var UrlSiteLocalHost = "localhost:4200"
var UrlSiteLocalHostIP = "127.0.0.1:4200"
var UrlSiteProduction = "54.156.244.197"
var UrlSiteHomol = "18.210.18.180"
var UrlSiteExterno = "https://flytura.com"
var UrlSiteSubDominio = "https://app.flytura.com/"
var BucketProd = "flytura-bucket"
var BucketHomol = "teste"
var BucketRegion = "us-east-1"
var ImagesInvoices = "invoices"
var AKA = os.Getenv("AKID")
var SKA = os.Getenv("ASK")
var FileAwsS3URL = "https://flytura-bucket.s3.us-east-1.amazonaws.com"
var ConectionString = "mongodb://admin:secret@localhost:27017"
var DBName = "flytura"
var Environment = "Prod"
var UserDBTableName = "user"
var PurcharseRecordTableName = "purcharseRecord"
var TokenAccessTableName = "tokenAccess"
var StatusImportTableName = "statusImport"
var PerfilTableName = "perfil"
var ImagesDBTableName = "imagesDB"
var AirlineTableName = "airline"
var OutPutInvoicesTableName = "outPutInvoices"
var ConciliationTableName = "conciliation"
var LogsTableName = "logs"
var Fuso1 string = "Europe/Lisbon"
var Fuso2 string = "America/Mexico_City"

// Abaixo número máximo de colunas permitidas na planilha purcharseRecord
var NumberMaxColunsSheet = 3

var NumberMaxColunsSheetOutPutInvoices = 19

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 17/12/2025 11:50
Data Final da criação :  10/12/2025 12:03
*/
// diffHours retorna a diferença (Portugal - México) em horas inteiras
// para um instante específico (instante UTC). Use este valor para subtrações simples.
func DiffHours(instant time.Time, fuso1 string, fuso2 string) (int, error) {
	lis, err := time.LoadLocation("Europe/Lisbon")
	if err != nil {
		return 0, err
	}
	mex, err := time.LoadLocation("America/Mexico_City")
	if err != nil {
		return 0, err
	}

	_, offLis := instant.In(lis).Zone() // offset em segundos
	_, offMex := instant.In(mex).Zone()

	// diferença de horas: exemplo: Lisbon(0 ou +3600) - Mexico(-21600) = +6h ou +7h
	diff := (offLis - offMex) / 3600
	return diff, nil
}

// var ConectionString = "mongodb://localhost:27017"

func TokenValido(w http.ResponseWriter, r *http.Request) (bool, string) {

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		// http.Error(w, "Token não fornecido", httjp.StatusUnauthorized)
		return false, "Token não fornecido"
	}

	// O formato esperado do cabeçalho é "Bearer <token>"
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		// http.Error(w, "Token malformado", http.StatusUnauthorized)
		return false, "Token malformado"
	}

	// Validar o token JWT
	_, err := ValidateToken(tokenString)
	if err != nil {
		// http.Error(w, fmt.Sprintf("Token inválido: %v", err), http.StatusUnauthorized)
		return false, "Token inválido"
	}

	return true, "Token válido"
}

// var secretKey = []byte("my_secret_key") // Chave secreta para validar o JWT

// Função para verificar o token JWT
func ValidateToken(tokenString string) (*jwt.Token, error) {
	// Parsing e validação do token JWT
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Verifique se o método de assinatura é o correto (HMAC)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de assinatura inválido")
		}
		return SecretKey, nil
	})
	if err != nil {
		return nil, err
	}
	return token, nil
}

func FormataRetornoHTTP(w http.ResponseWriter, obj any, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[string]any{"message": obj})
}

func FormataRetornoHTTPGeneric(w http.ResponseWriter, bodyName string, body any, codHttp int) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codHttp) // Código 200 OK
	return json.NewEncoder(w).Encode(map[any]any{"users": body})
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 14/10/2025 22:24
Data Final da criação :  14/10/2025 22:32
*/
func VerifyAccessValidTokenListSheet(client *mongo.Client, dbName, collectionName, token string) (map[string]any, error) {

	collection := client.Database(dbName).Collection(collectionName)

	// objectID, erroId := primitive.ObjectIDFromHex(excelId)
	// if erroId != nil {
	// 	log.Fatalf("Erro ao converter string para ObjectID: %v", erroId)
	// }

	filter := bson.M{
		"token":  token,
		"active": true,
	}

	// Variável para armazenar o usuário retornado
	var tokenAccessSheet models.TokenAccessSheet

	// Usar FindOne para pegar apenas um único registro
	err := collection.FindOne(context.Background(), filter).Decode(&tokenAccessSheet)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("token não encontrado")
		}
		return nil, fmt.Errorf("erro ao buscar plano de contas: %v", err)
	}

	// Converter o _id para string

	// Retornar o usuário como um mapa
	_return := map[string]any{
		"ID":   tokenAccessSheet.ID.Hex(), // Agora o campo ID é uma string
		"Name": tokenAccessSheet.Name,
		"Code": tokenAccessSheet.Token,
	}

	return _return, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 28/11/2025 11:28
Data Final da criação :  28/11/2025 11:28
*/

func GerarTokenSemExpiracao() (string, error) {
	// Cria as claims (dados do token)
	claims := jwt.MapClaims{
		"authorized": true,
		"user":       "Flytura",
		"iat":        time.Now().Unix(), // Data de emissão
		// Não adicionamos "exp" para não expirar
	}

	// Cria o token com método HMAC
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Assina o token com a chave secreta
	tokenString, err := token.SignedString(SecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Definir um mapa com os dados
	data := map[string]string{
		"name":   "Flytura API's",
		"Versao": "1.0",
		"Status": "OK",
	}

	// Definir o header como JSON
	w.Header().Set("Content-Type", "application/json")

	// Codificar o mapa de dados para JSON e enviar como resposta
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

func GetPublicIP() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	ip, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(ip), nil
}

var UploadChan chan UploadTask

type UploadTask struct {
	FileContent []byte
	FileName    string
	CompanyCode string
	Key         string
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/12/2025 12:54
Data Final da criação :  10/12/2025 12:58
*/
func ReturnBucketEnvironment() string {
	if Environment == "Prod" {
		return BucketProd
	} else {
		return BucketHomol
	}
}

// Classe de caracteres para um rune com suporte a acentos (pt/es) e case-agnóstico via Option "i".
func runeClassAcentoAgnostico(r rune) string {
	// Mapa simplificado; amplie conforme seu dataset
	switch unicode.ToLower(r) {
	case 'a':
		return "[aàáâãäåAÀÁÂÃÄÅ]"
	case 'e':
		return "[eèéêëEÈÉÊË]"
	case 'i':
		return "[iìíîïIÌÍÎÏ]"
	case 'o':
		return "[oòóôõöOÒÓÔÕÖ]"
	case 'u':
		return "[uùúûüUÙÚÛÜ]"
	case 'c':
		return "[cçCÇ]"
	case 'n':
		return "[nñNÑ]"
	default:
		// Escapa metacaracteres para não quebrar o regex
		return regexp.QuoteMeta(string(r))
	}
}

// Gera um padrão regex que ignora acentos E ignora espaços/separadores entre caracteres/grupos.
// - Colapsa qualquer sequência de espaços do termo em um único `\s*` no padrão.
// - Trata também separadores comuns ('-', '_', '.', '/', '\') como “espaço” opcional.
// - Use com Options: "i" para case-insensitive.
func AccentAgnosticSpaceInsensitivePattern(s string) string {
	// Quais caracteres do termo serão tratados como “espaço” (opcional) no alvo.
	sepAsSpace := map[rune]bool{
		' ': true, '\t': true, '\n': true, '\r': true, '\f': true,
		'-': true, '_': true, '.': true, '/': true, '\\': true,
	}

	var b strings.Builder
	addGap := false // indica se devemos inserir um \s* antes do próximo token não-separador

	for _, r := range s {
		if sepAsSpace[r] || unicode.IsSpace(r) {
			// Marca que precisamos admitir um gap antes do próximo “token”
			addGap = true
			continue
		}

		// Se havia espaço(s)/separador(es) antes, permita `\s*` aqui
		if addGap && b.Len() > 0 {
			b.WriteString(`\s*`)
			addGap = false
		}

		// Adiciona a classe acento-agnóstica para o rune atual
		b.WriteString(runeClassAcentoAgnostico(r))
	}

	// Se a string terminar com espaços/separadores, não precisamos fazer nada extra.
	// O padrão resultante permite \s* apenas entre tokens.

	return b.String()
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 26/11/2025 13:47
Data Final da criação :  26/11/2025 13:47
*/

func CriarArquivoTemporarioExcel(extensao string, nameFile string) string {
	if extensao == ".xls" {
		return nameFile + "-*.xls"
	}
	return nameFile + "-*.xlsx"
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 22/05/2026 16:45
Data Final da criação :  22/05/2026 16:45
*/
func Normalize(input string) string {
	// 1. Decomposição (separa acento da letra)
	t := norm.NFD.String(input)

	// 2. Remove os acentos
	var b strings.Builder
	for _, r := range t {
		if unicode.Is(unicode.Mn, r) {
			continue
		}
		b.WriteRune(r)
	}

	// 3. Minúsculo + remove espaços
	result := strings.ToLower(b.String())
	result = strings.ReplaceAll(result, " ", "")

	return result
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 01/03/2023 19:05
Data Final da criação : 01/03/2023 19:06
Local: Brasil
*/
func GetXlsColSafe(r *xls.Row, idx int) string {
	if r == nil {
		return ""
	}
	// LastCol geralmente é a contagem de colunas válidas (0..LastCol-1)
	if idx >= 0 && idx < r.LastCol() {
		return r.Col(idx)
	}
	return ""
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 01/03/2023 19:05
Data Final da criação : 01/03/2023 19:06
Local: Brasil
*/
// safeCell retorna a célula idx de uma []string (excelize GetRows)
// Se a coluna não existir, devolve "".
func SafeCell(row []string, idx int) string {
	if idx >= 0 && idx < len(row) {
		return strings.TrimSpace(row[idx])
	}
	return ""
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 01/03/2023 19:05
Data Final da criação : 01/03/2023 19:06
Local: Brasil
*/
func CompressSpaces(s string) string {

	// compressSpaces remove espaços repetidos no meio da string.
	var reSpaces = regexp.MustCompile(`\s+`)
	s = strings.TrimSpace(s)
	if s == "" {
		return s
	}
	return reSpaces.ReplaceAllString(s, "")
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 18:00
Data Final da criação 10/05/2026 18:00
*/
func ConvertToString(row map[string]bigquery.Value, fieldName string) string {

	result := ""
	if v, ok := row[fieldName].(string); ok {
		result = v
	}

	return result
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 18:00
Data Final da criação 10/05/2026 18:00
*/
func ConvertToTimeDate(row map[string]bigquery.Value, fieldName string) time.Time {

	if d, ok := row[fieldName].(civil.Date); ok {
		t := time.Date(
			d.Year,
			time.Month(d.Month),
			d.Day,
			0, 0, 0, 0,
			time.UTC,
		)

		return t
	} else {
		return time.Time{}
	}

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/05/2026 19:40
Data Final da criação 10/05/2026 19:40
*/
func ConvertToDecimal128(
	row map[string]bigquery.Value,
	fieldName string,
) primitive.Decimal128 {

	v, exists := row[fieldName]
	if !exists || v == nil {
		return primitive.NewDecimal128(0, 0)
	}

	switch value := v.(type) {

	// ✅ BigQuery NUMERIC / BIGNUMERIC
	case *big.Rat:
		d128, err := primitive.ParseDecimal128(
			value.FloatString(2),
		)
		if err != nil {
			return primitive.NewDecimal128(0, 0)
		}
		return d128

	// ✅ JSON / Excel como string ("3402.66" ou "3402,66")
	case string:
		value = strings.TrimSpace(value)
		if value == "" {
			return primitive.NewDecimal128(0, 0)
		}

		value = strings.Replace(value, ",", ".", 1)

		d128, err := primitive.ParseDecimal128(value)
		if err != nil {
			return primitive.NewDecimal128(0, 0)
		}
		return d128

	// ✅ Caso venha como número (menos comum, mas acontece)
	case float64:
		d128, err := primitive.ParseDecimal128(
			strconv.FormatFloat(value, 'f', -1, 64),
		)
		if err != nil {
			return primitive.NewDecimal128(0, 0)
		}
		return d128

	default:
		// tipo inesperado → não quebra o fluxo
		return primitive.NewDecimal128(0, 0)
	}
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 12/05/2026 18:47
Data Final da criação 12/05/2026 18:47
*/
func ConvertToFloat64(row map[string]bigquery.Value, fieldName string) float64 {
	if v, ok := row[fieldName]; ok {
		switch value := v.(type) {

		case float64:
			return value

		case int64:
			return float64(value)

		case string:
			// remove possíveis separadores de milhar
			clean := strings.ReplaceAll(value, ",", "")
			f, err := strconv.ParseFloat(clean, 64)
			if err == nil {
				return f
			}
		}
	}

	return 0
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/06/2026 10:34
Data Final da criação 10/06/2026 10:35
*/
func ConvertStringToTimeDate(value string) (time.Time, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return time.Time{}, errors.New("data vazia")
	}

	t, err := time.Parse("2006-01-02", value)
	if err != nil {
		return time.Time{}, err
	}

	return t, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/06/2026 10:34
Data Final da criação 10/06/2026 10:34
*/
func ConvertToInt(value string) (int, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return 0, errors.New("valor vazio")
	}

	i, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	return i, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/06/2026 10:39
Data Final da criação 10/06/2026 10:39
*/
func ConvertStringToFloat64(value string) (float64, error) {
	value = strings.TrimSpace(value)

	if value == "" {
		return 0, errors.New("valor vazio")
	}

	// remove separador de milhar (ex: 1,234.56)
	clean := strings.ReplaceAll(value, ",", "")

	f, err := strconv.ParseFloat(clean, 64)
	if err != nil {
		return 0, err
	}

	return f, nil
}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/06/2026 13:18
Data Final da criação 10/06/2026 13:18
*/
func ConvertNumberToDate(value string) time.Time {

	f, err := strconv.ParseFloat(value, 64)

	if err != nil {
		if value != "" && CompressSpaces(value) != "" {
			dt, err := ConvertStringToTimeDate(CompressSpaces(value)) // pode ficar vazio sem erro
			if err != nil {
				return time.Time{}
			} else {
				return dt
			}
		}
	} else {
		t, err := excelize.ExcelDateToTime(f, false)
		if err != nil {
			if value != "" && CompressSpaces(value) != "" {
				dt, err := ConvertStringToTimeDate(CompressSpaces(value)) // pode ficar vazio sem erro
				if err != nil {
					return time.Time{}
				} else {
					return dt
				}
			}
		} else {
			return t
		}

	}

	return time.Time{}

}

/*
Função criada por Ricardo Silva Ferreira
Inicio da criação 10/06/2026 13:22
Data Final da criação 10/06/2026 13:25
*/

func ConvertMonthToNumber(month string) (int, error) {
	month = strings.TrimSpace(strings.ToLower(month))

	months := map[string]int{
		"janeiro": 1, "jan": 1,
		"fevereiro": 2, "fev": 2,
		"março": 3, "mar": 3,
		"marco": 3, // sem acento
		"abril": 4, "abr": 4,
		"maio": 5, "mai": 5,
		"junho": 6, "jun": 6,
		"julho": 7, "jul": 7,
		"agosto": 8, "ago": 8,
		"setembro": 9, "set": 9,
		"outubro": 10, "out": 10,
		"novembro": 11, "nov": 11,
		"dezembro": 12, "dez": 12,
	}

	if val, ok := months[month]; ok {
		return val, nil
	}

	return 0, errors.New("mês inválido")
}
