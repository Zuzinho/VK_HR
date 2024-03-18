package handler

import (
	"VK_HR/pkg/userrepo"
	"VK_HR/pkg/validator"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func httpError(w http.ResponseWriter, err error, status int) {
	http.Error(w, fmt.Sprintf("{\"err\": \"%s\"}", err.Error()), status)
}

func write(w http.ResponseWriter, value any) error {
	buf, err := json.Marshal(value)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(buf)

	return err
}

func read(r *http.Request, form any) error {
	buf, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		return err
	}

	return json.Unmarshal(buf, form)
}

func (*AppHandler) checkAuth(r *http.Request) error {
	role, ok := r.Context().Value("user_role").(userrepo.Role)
	if !ok {
		return NoUserRoleErr
	}

	return role.IsValid()
}

func (*AppHandler) getIDFromQuery(queryParams map[string][]string, key string) (int, error) {
	idStr, exist := queryParams[key]
	if !exist {
		return 0, newNoRequiredParamError(key)
	}

	id, err := strconv.Atoi(idStr[0])
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (*AppHandler) getUpdateQueryParams(updateParams map[string][]string, valid validator.ValueValidator) (*string, *[]any, error) {
	var queryParts []string
	var args []interface{}
	i := 1

	var err error
	for column, values := range updateParams {
		value := values[0]
		var val any
		if val, err = valid.IsValidValue(validator.ColumnName(column), value); err != nil {
			return nil, nil, err
		}

		queryParts = append(queryParts, fmt.Sprintf("%s = $%d", column, i))
		args = append(args, val)
		i++
	}

	setClause := strings.Join(queryParts, ", ")

	return &setClause, &args, nil
}
