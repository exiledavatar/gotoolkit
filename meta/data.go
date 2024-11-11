package meta

type Data []any

func ToData(value any) Data {
	// var data Data
	// for _, v := range value {
	// 	data = append(data, v)
	// }

	// return data
	if value == nil {
		return nil
	}
	// fmt.Printf("Value is %#v -----------------------------------\n", value)

	return Data(ToSlice(value))
	// return Data{}
}
