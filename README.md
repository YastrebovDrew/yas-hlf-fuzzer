go install /go-fuzz-hf/gen  /n
go install /go-fuzz-hf/manifestgen /n
go install /go-fuzz-hf/go-fuzz /n
go install /go-fuzz-hf/go-fuzz-build /n
/n
cd ~/atb/fuzz /n
manifestgen -out=.  /n
*Исправляем манифест под наши нужды* /n
gen -manifest=manifest.json -out=corpus -limit=300  /n
go-fuzz-build-hf   -preserve=crypto/internal/bigmod  -o atb-fuzz.zip  /n 
go-fuzz-hf   -bin=atb-fuzz.zip     -procs=1  /n
