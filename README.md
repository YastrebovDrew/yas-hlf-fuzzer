Hyperledger Fabric Chaincode Fuzzing Toolkit

Инструкция по установке и использованию инструментов для фаззинга смарт-контрактов Hyperledger Fabric с помощью go-fuzz-hf.
Установка утилит
# Установить генератор сидов (manifestgen)
go install github.com/ваш/репозиторий/go-fuzz-hf/manifestgen

# Установить универсальный генератор корпуса (gen)
go install github.com/ваш/репозиторий/go-fuzz-hf/gen

# Установить кастомный билд-скрипт go-fuzz
go install github.com/ваш/репозиторий/go-fuzz-hf/go-fuzz

# Установить кастомный билд-пакер go-fuzz-build
go install github.com/ваш/репозиторий/go-fuzz-hf/go-fuzz-build

    Убедитесь, что $GOPATH/bin добавлен в ваш $PATH.

Подготовка к фаззингу

    Перейдите в директорию вашего chaincode‑харнесса (где находится fuzz.go):
    cd ~/hfuzz/atb/fuzz

    Создайте manifest.json с примером функций и их аргументов:
    manifestgen -out=.

    Отредактируйте manifest.json под ваш контракт, указав все публичные транзакции и диапазоны значений аргументов.

    Сгенерируйте сиды (corpus):
    gen -manifest=manifest.json -out=corpus -limit=300

    Здесь:

        -manifest — путь к манифесту;

        -out — директория для корпуса (будет создана, если не существует);

        -limit — максимум сидов на каждую функцию (по умолчанию 200).

Сборка и запуск go-fuzz

    Соберите бинарник фузз‑теста:
    go-fuzz-build -preserve=crypto/internal/bigmod -o atb-fuzz.zip

    Опция -preserve необходима для исключения конфликтов символов в internal.

    Запустите фуззер:
    go-fuzz -bin=atb-fuzz.zip -procs=1

        -procs=1 отключает параллельные гонки и защищает внутренний sonar от ошибок.

Теперь go-fuzz начнёт мутацию сидов из corpus и поиск аномалий в вашем смарт-контракте. Удачного фаззинга!

