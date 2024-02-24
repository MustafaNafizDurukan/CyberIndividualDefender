package chrome

import (
	"crypto/aes"
	"crypto/cipher"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"os"
	"path/filepath"

	"github.com/MustafaNafizDurukan/CyberIndividualDefender/pkg/fsutils"

	_ "github.com/mattn/go-sqlite3"
	"github.com/zavla/dpapi"
)

type Chrome struct {
	encKey         []byte
	LocalStatePath string
	LoginDataPath  string
}

func NewWithPath(localStatePath string, userDataPath string) *Chrome {
	return &Chrome{
		LocalStatePath: localStatePath,
		LoginDataPath:  userDataPath,
	}
}

func New() *Chrome {
	basePath := filepath.Join(os.Getenv("localappdata"), "Google", "Chrome", "User Data")

	return &Chrome{
		LocalStatePath: filepath.Join(basePath, "Local State"),
		LoginDataPath:  filepath.Join(basePath, "Default", "Login Data"),
	}
}

func (c *Chrome) Init() error {
	err := c.encryptionKey()
	if err != nil {
		return err
	}

	return nil
}

func (c *Chrome) ChromePasswords() ([]PasswordModel, error) {
	if dir, err := os.Stat(c.LoginDataPath); err != nil || !dir.IsDir() {
	}

	tempFile, _ := os.CreateTemp(os.Getenv("temp"), "")
	if err := fsutils.CopyFile(c.LoginDataPath, tempFile.Name()); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", tempFile.Name())
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query("SELECT origin_url, username_value, password_value FROM logins")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	passwords := []PasswordModel{}
	for rows.Next() {
		var url, username string
		var password []byte

		if err := rows.Scan(&url, &username, &password); err != nil {
			continue
		}

		password, err = c.DecryptPassword(password)
		if err != nil {
			continue
		}

		if url != "" && username != "" && len(password) != 0 {
			passwords = append(passwords, PasswordModel{
				URL:      url,
				Username: username,
				Password: string(password),
			})
		}
	}

	return passwords, nil
}

func (c *Chrome) DecryptPassword(encryptedData []byte) ([]byte, error) {
	if len(encryptedData) < 16 {
		return nil, errors.New("encrypted data is too short")
	}

	block, err := aes.NewCipher(c.encKey)
	if err != nil {
		return nil, err
	}

	GCMcipher, err := cipher.NewGCM(block)
	if err != nil {
		return nil, nil
	}

	decryptedData, err := GCMcipher.Open(nil, encryptedData[3:15], encryptedData[15:], nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func (c *Chrome) encryptionKey() error {
	if len(c.encKey) != 0 {
		return nil
	}

	file, err := os.Open(c.LocalStatePath)
	if err != nil {
		return err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var conf map[string]any
	if err := json.Unmarshal(content, &conf); err != nil {
		return err
	}

	v, ok := conf["os_crypt"]
	if !ok {
		return errors.New("chrome encryption key not found (os_crypt)")
	}

	osCrypt, ok := v.(map[string]any)
	if !ok {
		return errors.New("could not cast os_crypt")
	}

	eVal, ok := osCrypt["encrypted_key"]
	if !ok {
		return errors.New("key `encrypted_key` not found")
	}

	encKey, ok := eVal.(string)
	if !ok {
		return errors.New("could not cast encyption key into string")
	}

	encryptionKey, err := base64.StdEncoding.DecodeString(encKey)
	if err != nil {
		return err
	}

	encryptionKey, err = dpapi.Decrypt(encryptionKey[5:])
	if err != nil {
		return err
	}

	c.encKey = encryptionKey
	return nil
}
