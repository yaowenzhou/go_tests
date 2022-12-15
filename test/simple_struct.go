package test

type SimpleStruct struct {
	A    int    `json:"a" search:"lt,lte,eq,gte,gt,neq"`
	B    int    `json:"b" search:"lt,lte,eq,gte,gt,neq"`
	Str  string `json:"str" search:"contain,notcontain"`
	StrA string `json:"str_a"`
}
