package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

type FolderFile struct {
	Id         int    `orm:"column(file_id);auto" description:"文件id"`
	FileName   string `orm:"column(file_name);size(255)" description:"文件名称"`
	FileType   string `orm:"column(file_type);size(20)" description:"文件类型"`
	FileStatus string `orm:"column(file_status);size(20)" description:"文件状态"`
	FilePath   string `orm:"column(file_path);size(255)" description:"文件位置"`
	IdnumHash  string `orm:"column(idnum_hash);size(32)" description:"文件MD5 hash"`
	AddTime    int64  `orm:"column(add_time)" description:"添加时间"`
	UpdateTime int64  `orm:"column(update_time)" description:"更新时间"`
}

func (t *FolderFile) TableName() string {
	return "folder_file"
}

func init() {
	orm.RegisterModel(new(FolderFile))
}

// AddFolderFile insert a new FolderFile into database and returns
// last inserted Id on success.
func AddFolderFile(m *FolderFile) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetFolderFileById retrieves FolderFile by Id. Returns error if
// Id doesn't exist
func GetFolderFileById(id int) (v *FolderFile, err error) {
	o := orm.NewOrm()
	v = &FolderFile{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllFolderFile retrieves all FolderFile matches certain condition. Returns empty list if
// no records exist
func GetAllFolderFile(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(FolderFile))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []FolderFile
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateFolderFile updates FolderFile by Id and returns error if
// the record to be updated doesn't exist
func UpdateFolderFileById(m *FolderFile) (err error) {
	o := orm.NewOrm()
	v := FolderFile{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteFolderFile deletes FolderFile by Id and returns error if
// the record to be deleted doesn't exist
func DeleteFolderFile(id int) (err error) {
	o := orm.NewOrm()
	v := FolderFile{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&FolderFile{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
