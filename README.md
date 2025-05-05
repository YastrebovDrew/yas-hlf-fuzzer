
Инструкция по установке и использованию инструментов для фаззинга смарт-контрактов Hyperledger Fabric
# Переходим в папку фаззера
cd ~/yas-hlf-fuzzer/go-fuzz-hf
# Создать бинарник для генератора манифеста (manifestgen)
go install ./gen

# Создать бинарник для генератора начального корпуса (gen)
go install ./gen

# Создать бинарник go-fuzz
go install ./go-fuzz

# Создать бинарник билд-пакера go-fuzz-build
go install ./go-fuzz-build

    Убедитесь, что $GOPATH/bin добавлен в ваш $PATH.

Подготовка к фаззингу

    Перейдите в директорию вашего сhaincode‑харнесса (где находится fuzz.go):
    cd ~/yas-hlf-fuzzer/atb/fuzz

    Создайте manifest.json:
    manifestgen -out=.

    Отредактируйте manifest.json под ваш контракт, указав все функции и диапазоны значений аргументов.

    Сгенерируйте сиды (corpus):
    gen -manifest=manifest.json -out=corpus -limit=300 -limit-per-func=40

    Здесь:

        -manifest — путь к манифесту;

        -out — директория для корпуса (будет создана, если не существует);

        -limit — максимум сидов на каждую функцию (по умолчанию 200).

        -limit-per-func — максимум сидов на одну функцию

Сборка и запуск go-fuzz

    Соберите бинарник фузз‑теста:
    go-fuzz-build -tags=purego  -o chaincode-fuzz.zip

        Опция -tags=purego заставит компилятор игнорировать nat_asm.go и собрать «чисто-Go» реализацию. 

    Запустите фуззер:
    go-fuzz -bin=atb-fuzz.zip -procs=4


