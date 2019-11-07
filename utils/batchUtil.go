package utils

import (
	"bytes"
	"errors"
	"github.com/spf13/cast"
)

func NewInsertBuilder() *batchBuild {
	return &batchBuild{}
}

type batchBuild struct {
	Sql   bytes.Buffer //sql字符串
	Len   int          //插入列数
	Flag  bool
	Err   error
	Count int
}

func (s *batchBuild) Table(Table string) {
	s.Sql.WriteString("insert into " + Table)
}
func (s *batchBuild) Columns(Columns ...string) {
	if len(Columns) == 0 {
		s.Err = errors.New("插入列不能为空")
		return
	}
	s.Len = len(Columns)
	s.Sql.WriteString("( ")
	for i := range Columns {
		s.Sql.WriteString(Columns[i])
		if i != len(Columns)-1 {
			s.Sql.WriteString(",")
		} else {
			s.Sql.WriteString(")")
		}
	}
}
func (s *batchBuild) Values(Values ...interface{}) {
	if len(Values) == 0 {
		s.Err = errors.New("插入值不能为空")
		return
	}
	if s.Count != 0 {
		s.Sql.WriteString(",")
	} else {
		s.Sql.WriteString(" Values")
	}
	s.Count++
	if len(Values) != s.Len && s.Len != 0 {
		s.Err = errors.New("插入值与列数不匹配")
		return
	}
	s.Sql.WriteString("(")
	if s.Len == 0 {
		s.Sql.WriteString("null,Now(),null,null,")
	}
	for i := range Values {
		switch Values[i].(type) {
		case string:
			s.Sql.WriteString("'" + cast.ToString(Values[i]) + "'")
		default:
			s.Sql.WriteString(cast.ToString(Values[i]))
		}

		if i != len(Values)-1 {
			s.Sql.WriteString(",")
		} else {
			s.Sql.WriteString(")")
		}
	}
	s.Flag = true
}

/*
batchbuild := utils.NewInsertBuilder()
	batchbuild.Table("bills")
	batchbuild.Columns("created_at", "in_case_number", "bill_description", "case_no", "bill_source", "client_company_id", "accident_spot_description", "accident_spot_coordinate",
		"destination_coordinate", "destination_description", "customer_name", "customer_cellphone", "customer_license_plate", "customer_vehicle_type", "stats")
	for _, v := range bills {

		batchbuild.Values(time.Now().Format("2006-01-02 15:04:05"), v.InCaseNumber, v.BillDescription, v.CaseNo, v.BillSource, v.ClientCompanyID, v.AccidentSpotDescription, v.AccidentSpotCoordinate, v.DestinationCoordinate,
			v.DestinationDescription, v.CustomerName, v.CustomerCellphone, v.CustomerLicensePlate, v.CustomerVehicleType, 0)
	}
	if batchbuild.Flag {
		err := repo.GetDB().Exec(batchbuild.Sql.String()).Error
		if err != nil {
			fmt.Println(err)
		}
	}else{
		batchbuild.Err
	}
*/
