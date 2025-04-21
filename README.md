go install /go-fuzz-hf/gen
go install /go-fuzz-hf/manifestgen
go install /go-fuzz-hf/go-fuzz
go install /go-fuzz-hf/go-fuzz-build

cd ~/atb/fuzz
manifestgen -out=.
*Исправляем манифест под наши нужды*
gen -manifest=manifest.json -out=corpus -limit=300
go-fuzz-build-hf   -preserve=crypto/internal/bigmod  -o atb-fuzz.zip 
go-fuzz-hf   -bin=atb-fuzz.zip     -procs=1
