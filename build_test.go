package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

const (
	historicalV6SnippetSHA256           = "5840e7cc9ff1e1ba17413f2c54f242224aac413a94afb89a77cc6cf89c668d4b"
	historicalV6CommandsSHA256          = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration1SHA256          = "42f7919542f6ffb7d48fb8e44fd3a4091eb4a63f3d9774b594d526f56b6408b1"
	approvedV7Iteration1Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration2SHA256          = "7ddfd8260bdc5bddd951abe4712582b6e164d6ecdab693fcee7275ae526827b8"
	approvedV7Iteration2Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration3SHA256          = "04a2f6f0d046a828ac0d48d94e08469184232879e75859140ab28fbb28130902"
	approvedV7Iteration3Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration4SHA256          = "5929c529f11a3e9c01c7d80289025873ce9b8613e7db3d20606d0d8f70c15e32"
	approvedV7Iteration4Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration5SHA256          = "3310be931efc6aba4a25804391a5011c956f698ce0818d2a36d115bcb1f930c7"
	approvedV7Iteration5Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration6SHA256          = "fdbb2248248bd8dadf5108b82a8521b34bb7ae7c0fcffc0a81e2ec18a5f7fbf9"
	approvedV7Iteration6Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration7SHA256          = "d64627c3e4250f0e26bcf15a97b20c4a75a98dec3eb54ecb87b6d461de5bc1f5"
	approvedV7Iteration7Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration8SHA256          = "ce523b8dfe3e1d4c8c9f99ad93ea7d9d788bae0248bfff4aeb12706fbeb61c33"
	approvedV7Iteration8Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration9SHA256          = "7d09714ae18decf9ca6bbfb362049a78db77484e94133a73926c6e4e4e524dc6"
	approvedV7Iteration9Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration10SHA256         = "4b6863dc69dd9718876c511939ad53fb733e5ba099290c925e957199919203a4"
	approvedV7Iteration10Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration11SHA256         = "3dc4d44814da2e055820228717292ba482c0d65b0e612edefceb780e918daedf"
	approvedV7Iteration11Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration12SHA256         = "a0dbbe3cbca7fdc70560af96cb418a9d31ba335b6d5bf84537917d8affeaa9f7"
	approvedV7Iteration12Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration13SHA256         = "96d55eb656e9e48df559c38754be78396a0580e75653be4088300e36b5a85159"
	approvedV7Iteration13Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration14SHA256         = "0a3893084a3eb00d4ba389be1a156504191120ecea0cd28096a78e9af4e40145"
	approvedV7Iteration14Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration15SHA256         = "6b984c4bdd26d09e5b983b0d268cc8c11f32223683c9cfcbde5b414dbccf8893"
	approvedV7Iteration15Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration16SHA256         = "bcf18df13a9d6feec94697eaab51d65741cad95ce6227381a5bd8fbbdfc3ffa1"
	approvedV7Iteration16Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration17SHA256         = "3ed152152d9c755d2c18f0df4e1928d838812a6a04a33aa1ba5d9c37d8c3e095"
	approvedV7Iteration17Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration18SHA256         = "602bf3edd6dfef153dfcaf0166ca3f0e56927b951102ea59f9a67088ff77eb11"
	approvedV7Iteration18Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration19SHA256         = "e60bbdf48c50df882774e7b08402c9aa11346880ab65561b0d23271e4943e305"
	approvedV7Iteration19Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedV7Iteration20SHA256         = "eb329317f710bc17689d92ef768348401c6bfa135dd2966c18dd00780c855c9c"
	approvedV7Iteration20Commands       = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration21SHA256          = "46ea401b4f979ec857b2a079308c426f2f95fa0cc54c3f6201ec162ca398d111"
	currentV7Iteration21Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration22SHA256          = "303cc7ea8f7d18cc31408f7bdf3424b31711fdbc5f68c890590941bb62c0f189"
	currentV7Iteration22Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration23SHA256          = "0ca1bd5d0dbaf01a2ef5371a88f8f59c01593296192a0f21ba0a4321e30236a2"
	currentV7Iteration23Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration24SHA256          = "85f91833817a9fb9fa422d6f395743f1397b9e8b3ab7fcf6283d1efcc4f817d4"
	currentV7Iteration24Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration25SHA256          = "d6ae6e6fd6bafcd949c643e0fa6edb740e8e7ca8a3451ac7f044bc71cc7a1674"
	currentV7Iteration25Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration26SHA256          = "c3168fabd9761906a6004c3a689827968caabb46607cb9815ac4190f52d718e2"
	currentV7Iteration26Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration27SHA256          = "5b5b8bc76050ef0e6725403cd3a7a3c8a2f3de03ea2f8f31b2641b5a59bbf0b0"
	currentV7Iteration27Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration28SHA256          = "743f9201b2058249222ec19a3136476b31c48e408c7a4d75bf85db220de58646"
	currentV7Iteration28Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration29SHA256          = "eaec72fbaa1334f789b5cc7451361b83bc920984fce6093ec88e38c85e384c7a"
	currentV7Iteration29Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration30SHA256          = "2c52360c9d0a68e0a028f02660b15242206f5cd2645484dcc00fad1334f1e03f"
	currentV7Iteration30Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration31SHA256          = "11488609a0ade7d3ea45a9c89ebc91c49270ce12b19c154528d700f73ac8d35a"
	currentV7Iteration31Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration32SHA256          = "a8640a242c972e62d900cee36fbe115a21b591730940af37f6b8a6b5888284ca"
	currentV7Iteration32Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration33SHA256          = "a7e0b39fa9ce92abdcec472d273db5ef8fd040080b10092727bd208e1f6dcdb8"
	currentV7Iteration33Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	currentV7Iteration34SHA256          = "3986dece5e7ddfca8cfffacc0007d977a8a70641aa817820b4537068eb6b8b6f"
	currentV7Iteration34Commands        = "d97aac2dfe4e8d21e29c892b2d24ba991f5b568f6f8225c724f357bc3efe0595"
	approvedMigrationName               = "archive.compress-tar-gz"
	iteration2MigrationName             = "archive.compress-tar-xz"
	iteration3MigrationName             = "archive.compress-zip"
	iteration4MigrationName             = "archive.decompress-tar-gz"
	iteration5MigrationName             = "archive.decompress-tar-xz"
	iteration6MigrationName             = "archive.decompress-unzip"
	iteration7MigrationName             = "array.all-elements"
	iteration8MigrationName             = "array.at-index"
	iteration9MigrationName             = "array.concat"
	iteration10MigrationName            = "array.contains"
	iteration11MigrationName            = "array.declare"
	iteration12MigrationName            = "array.delete-at"
	iteration13MigrationName            = "array.delete"
	iteration14MigrationName            = "array.filter"
	iteration15MigrationName            = "array.iterate"
	iteration16MigrationName            = "array.length"
	iteration17MigrationName            = "array.print"
	iteration18MigrationName            = "array.push"
	iteration19MigrationName            = "array.range"
	iteration20MigrationName            = "array.replace"
	iteration21MigrationName            = "array.reverse"
	iteration22MigrationName            = "array.set-element-at"
	iteration23MigrationName            = "command.failure-check"
	iteration24MigrationName            = "command.hide-error"
	iteration25MigrationName            = "command.if-exists"
	iteration26MigrationName            = "command.nice"
	iteration27MigrationName            = "command.renice"
	iteration28MigrationName            = "command.run"
	iteration29MigrationName            = "command.substitution"
	iteration30MigrationName            = "command.success-check"
	iteration31MigrationName            = "cryptography.base64-decode"
	iteration32MigrationName            = "cryptography.base64-encode"
	iteration33MigrationName            = "cryptography.hash"
	iteration34MigrationName            = "date.date-now-short"
	v6CompressTarGzBody                 = "tar -czvf ${1|/path/to/archive, \"${pathToArchive}\"|}.tar.gz ${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\n"
	currentCompressTarGzBody            = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nsource_path=\"${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\"\ncase ${source_path} in\n  -*) source_path=./${source_path} ;;\nesac\ntar -cf \"${archive_path}.tar\" \"${source_path}\" && gzip -f \"${archive_path}.tar\"\n"
	v6CompressTarXzBody                 = "tar -cJf ${1|/path/to/archive, \"${pathToArchive}\"|}.tar.xz ${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\n"
	currentCompressTarXzBody            = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nsource_path=\"${2|/path/to/directory-or-file, \"${pathToDirectoryOrFile}\"|}\"\ncase ${source_path} in\n  -*) source_path=./${source_path} ;;\nesac\ntar -cf \"${archive_path}.tar\" \"${source_path}\" && xz -f \"${archive_path}.tar\"\n"
	v6CompressZipBody                   = "zip -rq ${1|/path/to/archive, \"${pathToArchive}\"|}.zip ${2|/path/to/directory-or-file,\"${pathToDirectoryOrFile}\"|}\n"
	currentCompressZipBody              = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nsource_path=\"${2|/path/to/directory-or-file,\"${pathToDirectoryOrFile}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\ncase ${source_path} in\n  -*) source_path=./${source_path} ;;\nesac\nzip -rq \"${archive_path}.zip\" \"${source_path}\"\n"
	v6DecompressTarGzBody               = "tar -C ${1|/extract/to/path, \"${extractToPath}\"|} -xzvf ${2|/path/to/archive, \"${pathToArchive}\"|}.tar.gz\n"
	currentDecompressTarGzBody          = "extract_path=\"${1|/extract/to/path, \"${extractToPath}\"|}\"\narchive_path=\"${2|/path/to/archive, \"${pathToArchive}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\n(\n  temporary_directory=$(mktemp -d \"${TMPDIR:-/tmp}/shellman.XXXXXXXXXX\") || exit\n  status=0\n  trap 'status=$?; trap - 0; rm -rf \"${temporary_directory}\"; exit \"${status}\"' 0\n  trap 'exit 129' HUP\n  trap 'exit 130' INT\n  trap 'exit 143' TERM\n  gzip -dc \"${archive_path}.tar.gz\" > \"${temporary_directory}/archive.tar\" &&\n    (cd \"${extract_path}\" && tar -xf \"${temporary_directory}/archive.tar\")\n)\n"
	v6DecompressTarXzBody               = "tar -C ${1|/extract/to/path, \"${extractToPath}\"|} -xf ${2|/path/to/archive, \"${pathToArchive}\"|}.tar.xz\n"
	currentDecompressTarXzBody          = "extract_path=\"${1|/extract/to/path, \"${extractToPath}\"|}\"\narchive_path=\"${2|/path/to/archive, \"${pathToArchive}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\n(\n  temporary_directory=$(mktemp -d \"${TMPDIR:-/tmp}/shellman.XXXXXXXXXX\") || exit\n  status=0\n  trap 'status=$?; trap - 0; rm -rf \"${temporary_directory}\"; exit \"${status}\"' 0\n  trap 'exit 129' HUP\n  trap 'exit 130' INT\n  trap 'exit 143' TERM\n  xz -dc \"${archive_path}.tar.xz\" > \"${temporary_directory}/archive.tar\" &&\n    (cd \"${extract_path}\" && tar -xf \"${temporary_directory}/archive.tar\")\n)\n"
	v6DecompressUnzipBody               = "unzip -q ${1|/path/to/archive, \"${pathToArchive}\"|}.zip -d ${2|/extract/to/path,\"${extractToPath}\"|}"
	currentDecompressUnzipBody          = "archive_path=\"${1|/path/to/archive, \"${pathToArchive}\"|}\"\nextract_path=\"${2|/extract/to/path,\"${extractToPath}\"|}\"\ncase ${archive_path} in\n  -*) archive_path=./${archive_path} ;;\nesac\nunzip -q \"${archive_path}.zip\" -d \"${extract_path}\""
	v6ArrayAllElementsBody              = "${0:echo }\"${${1:myArray}[@]}\""
	currentArrayAllElementsBody         = "# Bash-only: indexed arrays are not specified by POSIX sh.\n${0:echo }\"${${1:myArray}[@]}\""
	v6ArrayAtIndexBody                  = "${0:echo }\"${${1:myArray}[${2:index}]}\""
	currentArrayAtIndexBody             = "# Bash-only: indexed arrays are not specified by POSIX sh.\n${0:echo }\"${${1:myArray}[${2:index}]}\""
	v6ArrayConcatBody                   = "${1:newArray}=(\"${${2:array1}[@]}\" \"${${3:array2}[@]}\")\n"
	currentArrayConcatBody              = "# Bash-only: indexed arrays are not specified by POSIX sh.\n${1:newArray}=(\"${${2:array1}[@]}\" \"${${3:array2}[@]}\")\n"
	v6ArrayDeleteAtBody                 = "unset \"${1:myArray}[${2:index}]\"\n"
	currentArrayDeleteAtBody            = "# Bash-only: indexed arrays are not specified by POSIX sh.\nunset \"${1:myArray}[${2:index}]\"\n"
	v6ArrayDeleteBody                   = "unset ${1:myArray}\n"
	currentArrayDeleteBody              = "# Bash-only: indexed arrays are not specified by POSIX sh.\nunset \"${1:myArray}\"\n"
	v6ArrayFilterBody                   = "readarray -t ${1:filtered} < <(for i in \"${${2:myArray}[@]}\" ; do echo \"\\${i\\}\"; done | grep ${3|',\"|}${4:pattern}${3})\n"
	currentArrayFilterBody              = "# Bash-only: indexed arrays, readarray, and process substitution are not specified by POSIX sh.\nreadarray -t ${1:filtered} < <(for i in \"${${2:myArray}[@]}\"; do printf '%s\\n' \"\\${i\\}\"; done | grep -e ${3|',\"|}${4:pattern}${3} || [ $? -eq 1 ]) && wait \"$!\"\n"
	v6ArrayLengthBody                   = "${1:length}=${#${2:myArray}[@]}\n"
	currentArrayLengthBody              = "# Bash-only: indexed arrays are not specified by POSIX sh.\n${1:length}=${#${2:myArray}[@]}\n"
	v6ArrayPrintBody                    = "echo \"\\${${1:myArray}[@]}\"\n"
	currentArrayPrintBody               = "# Bash-only: indexed arrays are not specified by POSIX sh.\n(IFS=' '; printf '%s\\n' \"\\${${1:myArray}[*]}\")\n"
	v6ArrayPushBody                     = "${1:myArray}+=('${2:newItem}')\n"
	currentArrayPushBody                = "# Bash-only: indexed arrays and compound assignment are not specified by POSIX sh.\n${1:myArray}+=('${2:newItem}')\n"
	v6ArrayRangeBody                    = "${1:newArray}=\"${${2:myArray}[*]:${3:fromIndex}:${4:n}}\"\n"
	currentArrayRangeBody               = "# Bash-only: indexed arrays and array slicing are not specified by POSIX sh.\n${1:newArray}=(\"${${2:myArray}[@]:${3:fromIndex}:${4:n}}\")\n"
	v6ArrayReplaceBody                  = "${1:newArray}=${${2:myArray}[*]//${3:find}/${4:replace}}\n"
	currentArrayReplaceBody             = "# Bash-only: indexed arrays and pattern substitution are not specified by POSIX sh.\n${1:newArray}=(\"${${2:myArray}[@]//${3:find}/${4:replace}}\")\n"
	v6ArraySetElementAtBody             = "${1:myArray}[${2:index}]=\"${3:value}\"\n"
	currentArraySetElementAtBody        = "# Bash-only: indexed arrays and arithmetic array indices are not specified by POSIX sh.\n${1:myArray}[${2:index}]=\"${3:value}\"\n"
	v6CommandHideErrorBody              = "${1:command} 2> /dev/null\n"
	currentCommandHideErrorBody         = "{\n\t${1:command}\n} 2>/dev/null\n"
	v6CommandNiceBody                   = "sudo nice -n ${1|-20,-15,-10,-5,0,5,10,15,19|} ${2:command}\n"
	currentCommandNiceBody              = "sudo nice -n ${1|-20,-15,-10,-5,0,5,10,15,19|} -- ${2:command}\n"
	v6CommandReniceBody                 = "for p in \\$(pidof \"${1:processName}\"); do sudo renice -n ${2|-20,-15,-10,-5,0,5,10,15,19|} -p \"\\$p\"; done\n"
	currentCommandReniceBody            = "(\n\t_shellman_renice_pids=\\$(pidof -- \"${1:processName}\") || exit\n\t_shellman_renice_status=0\n\tunset IFS\n\t# shellcheck disable=SC2086 # pidof output must split into individual PIDs\n\tfor _shellman_renice_pid in \\${_shellman_renice_pids}; do\n\t\tsudo renice -n ${2|-20,-15,-10,-5,0,5,10,15,19|} -p \"\\${_shellman_renice_pid}\" || _shellman_renice_status=\\$?\n\tdone\n\texit \"\\${_shellman_renice_status}\"\n)\n"
	v6CommandRunBody                    = "${1:result}=\"$(${2:command})\"\n"
	currentCommandRunBody               = "# POSIX sh: command substitution runs in a subshell and removes trailing newlines.\n${1:result}=\"$(${2:command})\"\n"
	v6CommandSubstitutionBody           = "${1:result}=\"$(${2:command})\"\n"
	currentCommandSubstitutionBody      = "# POSIX sh: command substitution runs in a subshell and removes trailing newlines.\n${1:result}=\"$(${2:command})\"\n"
	v6CryptographyBase64DecodeBody      = "${1:base64Decoded}=\\$(echo -n \"${2|stringToDecode,${variableToDecode}|}\" | base64 -d)\n"
	currentCryptographyBase64DecodeBody = "# Requires a base64 utility with the -d decode option.\n${1:base64Decoded}=\\$(printf '%s' \"${2|stringToDecode,${variableToDecode}|}\" | base64 -d)\n"
	v6CryptographyBase64EncodeBody      = "${1:base64Encoded}=\\$(echo -n \"${2|stringToEncode,${variableToEncode}|}\" | base64)\n"
	currentCryptographyBase64EncodeBody = "${1:base64Encoded}=\\$(printf '%s' \"${2|stringToEncode,${variableToEncode}|}\" | base64)\n"
	v6CryptographyHashBody              = "${1:hash}=\\$(echo -n \"\\$${2:variableToHash}\" | ${3|md5sum,shasum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum|} | cut -f1 -d ' ')\n"
	currentCryptographyHashBody         = "${1:hash}=\\$(\n  hash_output=$(printf '%s' \"\\$${2:variableToHash}\" | ${3|md5sum,shasum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum|}) || exit\n  printf '%s\\n' \"\\${hash_output%% *}\"\n)\n"
	v6DateNowShortBody                  = "${1:dateShort}=\\$(date -I) ${0:# format: yyyy/mm/dd}\n"
	currentDateNowShortBody             = "${1:dateShort}=\\$(date '+%Y/%m/%d') ${0:# format: yyyy/mm/dd}\n"
)

var v6ArrayReverseBody = []any{
	`for((i=\${#${1:myArray}[@]}-1;i>=0;i--)); do`,
	"\t${2:reversed}+=(\"\\${${1:myArray}[i]}\")",
	"done",
	"",
	"unset \"${1:myArray}\" # optional",
	"echo \"\\${${2:reversed}[@]}\"",
	"",
}

var currentArrayReverseBody = []any{
	"# Bash-only: indexed arrays and C-style for loops are not specified by POSIX sh.",
	"# shellcheck disable=SC2050 # compare expanded snippet variable names",
	"${2:reversed}=(\"\\${${1:myArray}[@]}\") &&",
	"_shellman_reverse_status=0 &&",
	`for ((i=0, j=\${#${2:reversed}[@]}-1; i<j; i++, j--)); do`,
	"\t_shellman_reverse_item=\\${${2:reversed}[i]} &&",
	"\t\t${2:reversed}[i]=\\${${2:reversed}[j]} &&",
	"\t\t${2:reversed}[j]=\\${_shellman_reverse_item} || { _shellman_reverse_status=$?; break; }",
	"done &&",
	"((! _shellman_reverse_status)) &&",
	"",
	"{ [ '${1:myArray}' = '${2:reversed}' ] || unset \"${1:myArray}\"; } && # optional",
	"(IFS=' '; printf '%s\\n' \"\\${${2:reversed}[*]}\")",
	"",
}

var v6CommandFailureCheckBody = []any{
	"if ! ${1:command} >/dev/null 2>&1; then",
	"\techo \"failed\"",
	"else",
	"\techo \"succeed\"",
	"fi\n",
}

var currentCommandFailureCheckBody = []any{
	"if ${1:command} >/dev/null 2>&1; then",
	"\t_shellman_command_status=0 &&",
	"\tprintf '%s\\n' 'succeed'",
	"else",
	"\t_shellman_command_status=$? &&",
	"\tprintf '%s\\n' 'failed'",
	"fi &&",
	"(exit \"\\${_shellman_command_status}\")\n",
}

var v6CommandSuccessCheckBody = []any{
	"if ${1:command} >/dev/null 2>&1; then",
	"\techo \"succeed\"",
	"else",
	"\techo \"failed\"",
	"fi\n",
}

var currentCommandSuccessCheckBody = []any{
	"if {",
	"\t${1:command}",
	"} >/dev/null 2>&1; then",
	"\techo \"succeed\"",
	"else",
	"\t(",
	"\t\tset -- \"\\$?\"",
	"\t\techo \"failed\" || exit",
	"\t\texit \"\\$1\"",
	"\t)",
	"fi\n",
}

var v6CommandIfExistsBody = []any{
	`if [ "\$(command -v ${1:command})" ]; then`,
	"\t${2:echo \"command \\\"${1:command}\\\" exists on system\"}",
	"fi\n",
}

var currentCommandIfExistsBody = []any{
	`if command -v -- "${1:command}" >/dev/null 2>&1; then`,
	"\t${2:echo \"command \\\"${1:command}\\\" exists on system\"}",
	"fi\n",
}

var v6ArrayContainsBody = []any{
	"if [[ \"\\${${1:myArray}[*]}\" =~ ${2|'element',\"${value}\"|} ]]; then",
	"\techo 'array contains element'",
	"fi\n",
}

var currentArrayContainsBody = []any{
	"# Bash-only: indexed arrays are not specified by POSIX sh.",
	"(",
	"\tfor element in \"${${1:myArray}[@]}\"; do",
	"\t\tif [ \"\\${element\\}\" = ${2|'element',\"${value}\"|} ]; then",
	"\t\t\techo 'array contains element'",
	"\t\t\texit",
	"\t\tfi",
	"\tdone",
	")\n",
}

var v6ArrayDeclareBody = []any{
	"${1:myArray}=(",
	"\t'${2:constant}'",
	"\t\"${3:${variable\\}}\"",
	"\t'${4:another constant}'",
	")\n",
}

var currentArrayDeclareBody = []any{
	"# Bash-only: indexed arrays are not specified by POSIX sh.",
	"${1:myArray}=(",
	"\t'${2:constant}'",
	"\t\"${3:${variable\\}}\"",
	"\t'${4:another constant}'",
	")\n",
}

var v6ArrayIterateBody = []any{
	`for item in "${${1:myArray}[@]}"; do`,
	"\t${2:echo \"\\${item\\}\"}",
	"done\n",
}

var currentArrayIterateBody = []any{
	"# Bash-only: indexed arrays are not specified by POSIX sh.",
	`for item in "${${1:myArray}[@]}"; do`,
	"\t${2:echo \"\\${item\\}\"}",
	"done\n",
}

var v6Namespaces = []string{
	"archive", "array", "command", "cryptography", "date", "event",
	"filesystem", "float", "fn-fx", "ftp", "function", "git", "http",
	"input", "integer", "internal", "ip", "math", "misc", "output",
	"process", "string", "system", "time", "variable",
}

func TestGeneratedOutputMatchesFiles(t *testing.T) {
	snippetJSON, commands, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}

	assertFileBytes(t, snippetOutputPath, snippetJSON)
	assertFileBytes(t, documentOutputPath, commands)
	if !bytes.HasSuffix(snippetJSON, []byte("\n")) {
		t.Error("snippet output must end with a newline")
	}
	if !bytes.HasSuffix(commands, []byte("\n\n")) {
		t.Error("documentation output must end with two newlines")
	}
}

func TestMigrationChangesOnlyApprovedSnippet(t *testing.T) {
	ordered, err := readSnippets(rootDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(ordered), 278; got != want {
		t.Fatalf("snippet count = %d, want %d", got, want)
	}
	currentOrdered := append([]namedSnippet(nil), ordered...)
	if got := ordered[3].name; got != iteration4MigrationName {
		t.Fatalf("fourth snippet in actual nsroot traversal = %q, want %q", got, iteration4MigrationName)
	}
	if got := ordered[4].name; got != iteration5MigrationName {
		t.Fatalf("fifth snippet in actual nsroot traversal = %q, want %q", got, iteration5MigrationName)
	}
	if got := ordered[5].name; got != iteration6MigrationName {
		t.Fatalf("sixth snippet in actual nsroot traversal = %q, want %q", got, iteration6MigrationName)
	}
	if got := ordered[6].name; got != iteration7MigrationName {
		t.Fatalf("seventh snippet in actual nsroot traversal = %q, want %q", got, iteration7MigrationName)
	}
	if got := ordered[7].name; got != iteration8MigrationName {
		t.Fatalf("eighth snippet in actual nsroot traversal = %q, want %q", got, iteration8MigrationName)
	}
	if got := ordered[8].name; got != iteration9MigrationName {
		t.Fatalf("ninth snippet in actual nsroot traversal = %q, want %q", got, iteration9MigrationName)
	}
	if got := ordered[9].name; got != iteration10MigrationName {
		t.Fatalf("tenth snippet in actual nsroot traversal = %q, want %q", got, iteration10MigrationName)
	}
	if got := ordered[10].name; got != iteration11MigrationName {
		t.Fatalf("eleventh snippet in actual nsroot traversal = %q, want %q", got, iteration11MigrationName)
	}
	if got := ordered[11].name; got != iteration12MigrationName {
		t.Fatalf("twelfth snippet in actual nsroot traversal = %q, want %q", got, iteration12MigrationName)
	}
	if got := ordered[12].name; got != iteration13MigrationName {
		t.Fatalf("thirteenth snippet in actual nsroot traversal = %q, want %q", got, iteration13MigrationName)
	}
	if got := ordered[13].name; got != iteration14MigrationName {
		t.Fatalf("fourteenth snippet in actual nsroot traversal = %q, want %q", got, iteration14MigrationName)
	}
	if got := ordered[14].name; got != iteration15MigrationName {
		t.Fatalf("fifteenth snippet in actual nsroot traversal = %q, want %q", got, iteration15MigrationName)
	}
	if got := ordered[15].name; got != iteration16MigrationName {
		t.Fatalf("sixteenth snippet in actual nsroot traversal = %q, want %q", got, iteration16MigrationName)
	}
	if got := ordered[16].name; got != iteration17MigrationName {
		t.Fatalf("seventeenth snippet in actual nsroot traversal = %q, want %q", got, iteration17MigrationName)
	}
	if got := ordered[17].name; got != iteration18MigrationName {
		t.Fatalf("eighteenth snippet in actual nsroot traversal = %q, want %q", got, iteration18MigrationName)
	}
	if got := ordered[18].name; got != iteration19MigrationName {
		t.Fatalf("nineteenth snippet in actual nsroot traversal = %q, want %q", got, iteration19MigrationName)
	}
	if got := ordered[19].name; got != iteration20MigrationName {
		t.Fatalf("twentieth snippet in actual nsroot traversal = %q, want %q", got, iteration20MigrationName)
	}
	if got := ordered[20].name; got != iteration21MigrationName {
		t.Fatalf("twenty-first snippet in actual nsroot traversal = %q, want %q", got, iteration21MigrationName)
	}
	if got := ordered[21].name; got != iteration22MigrationName {
		t.Fatalf("twenty-second snippet in actual nsroot traversal = %q, want %q", got, iteration22MigrationName)
	}
	if got := ordered[22].name; got != iteration23MigrationName {
		t.Fatalf("entry after array namespace = %q, want %q", got, iteration23MigrationName)
	}
	if got := ordered[23].name; got != iteration24MigrationName {
		t.Fatalf("twenty-fourth snippet in actual nsroot traversal = %q, want %q", got, iteration24MigrationName)
	}
	if got := ordered[24].name; got != iteration25MigrationName {
		t.Fatalf("twenty-fifth snippet in actual nsroot traversal = %q, want %q", got, iteration25MigrationName)
	}
	if got := ordered[25].name; got != iteration26MigrationName {
		t.Fatalf("twenty-sixth snippet in actual nsroot traversal = %q, want %q", got, iteration26MigrationName)
	}
	if got := ordered[26].name; got != iteration27MigrationName {
		t.Fatalf("twenty-seventh snippet in actual nsroot traversal = %q, want %q", got, iteration27MigrationName)
	}
	if got := ordered[27].name; got != iteration28MigrationName {
		t.Fatalf("twenty-eighth snippet in actual nsroot traversal = %q, want %q", got, iteration28MigrationName)
	}
	if got := ordered[28].name; got != iteration29MigrationName {
		t.Fatalf("twenty-ninth snippet in actual nsroot traversal = %q, want %q", got, iteration29MigrationName)
	}
	if got := ordered[29].name; got != iteration30MigrationName {
		t.Fatalf("thirtieth snippet in actual nsroot traversal = %q, want %q", got, iteration30MigrationName)
	}
	if got := ordered[30].name; got != iteration31MigrationName {
		t.Fatalf("thirty-first snippet in actual nsroot traversal = %q, want %q", got, iteration31MigrationName)
	}
	if got := ordered[31].name; got != iteration32MigrationName {
		t.Fatalf("thirty-second snippet in actual nsroot traversal = %q, want %q", got, iteration32MigrationName)
	}
	if got := ordered[32].name; got != iteration33MigrationName {
		t.Fatalf("thirty-third snippet in actual nsroot traversal = %q, want %q", got, iteration33MigrationName)
	}
	if got := ordered[33].name; got != iteration34MigrationName {
		t.Fatalf("thirty-fourth snippet in actual nsroot traversal = %q, want %q", got, iteration34MigrationName)
	}
	currentIteration34, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-34 snippets (thirty-four changed, 244 unchanged)", currentIteration34, currentV7Iteration34SHA256)
	assertSHA256(t, "current v7 iteration-34 commands", renderDocumentation(ordered), currentV7Iteration34Commands)

	foundFirst := false
	foundSecond := false
	foundThird := false
	foundFourth := false
	foundFifth := false
	foundSixth := false
	foundSeventh := false
	foundEighth := false
	foundNinth := false
	foundTenth := false
	foundEleventh := false
	foundTwelfth := false
	foundThirteenth := false
	foundFourteenth := false
	foundFifteenth := false
	foundSixteenth := false
	foundSeventeenth := false
	foundEighteenth := false
	foundNineteenth := false
	foundTwentieth := false
	foundTwentyFirst := false
	foundTwentySecond := false
	foundTwentyThird := false
	foundTwentyFourth := false
	foundTwentyFifth := false
	foundTwentySixth := false
	foundTwentySeventh := false
	foundTwentyEighth := false
	foundTwentyNinth := false
	foundThirtieth := false
	foundThirtyFirst := false
	foundThirtySecond := false
	foundThirtyThird := false
	foundThirtyFourth := false
	for index := range ordered {
		switch ordered[index].name {
		case approvedMigrationName:
			foundFirst = true
			wantPrefix := []any{"archive compress tar.gz", "archive tar.gz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("approved snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "compress file/folder to a .tar.gz file"; got != want {
				t.Fatalf("approved snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCompressTarGzBody {
				t.Fatalf("approved snippet body = %#v, want %#v", got, currentCompressTarGzBody)
			}
		case iteration2MigrationName:
			foundSecond = true
			wantPrefix := []any{"archive compress tar.xz", "archive tar.xz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-2 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "compress file/folder to a .tar.xz file"; got != want {
				t.Fatalf("iteration-2 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCompressTarXzBody {
				t.Fatalf("iteration-2 snippet body = %#v, want %#v", got, currentCompressTarXzBody)
			}
		case iteration3MigrationName:
			foundThird = true
			wantPrefix := []any{"archive compress .zip", "archive zip"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-3 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "compress file/folder to a .zip file"; got != want {
				t.Fatalf("iteration-3 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCompressZipBody {
				t.Fatalf("iteration-3 snippet body = %#v, want %#v", got, currentCompressZipBody)
			}
		case iteration4MigrationName:
			foundFourth = true
			wantPrefix := []any{"archive decompress tar.gz", "decompress tar.gz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-4 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "decompress a .tar.gz file to specified path"; got != want {
				t.Fatalf("iteration-4 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDecompressTarGzBody {
				t.Fatalf("iteration-4 snippet body = %#v, want %#v", got, currentDecompressTarGzBody)
			}
		case iteration5MigrationName:
			foundFifth = true
			wantPrefix := []any{"archive decompress tar.xz", "decompress tar.xz"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-5 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "decompress a .tar.xz file to specified path"; got != want {
				t.Fatalf("iteration-5 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDecompressTarXzBody {
				t.Fatalf("iteration-5 snippet body = %#v, want %#v", got, currentDecompressTarXzBody)
			}
		case iteration6MigrationName:
			foundSixth = true
			wantPrefix := []any{"archive decompress .zip", "archive unzip"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-6 snippet prefixes changed: got %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "decompress a .zip file to specified path"; got != want {
				t.Fatalf("iteration-6 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDecompressUnzipBody {
				t.Fatalf("iteration-6 snippet body = %#v, want %#v", got, currentDecompressUnzipBody)
			}
		case iteration7MigrationName:
			foundSeventh = true
			if got, want := ordered[index].snippet.Prefix, "array all"; got != want {
				t.Fatalf("iteration-7 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "access all array elements"; got != want {
				t.Fatalf("iteration-7 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayAllElementsBody {
				t.Fatalf("iteration-7 snippet body = %#v, want %#v", got, currentArrayAllElementsBody)
			}
		case iteration8MigrationName:
			foundEighth = true
			if got, want := ordered[index].snippet.Prefix, "array at index"; got != want {
				t.Fatalf("iteration-8 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "retrieve element from array at specified index (zero based)"; got != want {
				t.Fatalf("iteration-8 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayAtIndexBody {
				t.Fatalf("iteration-8 snippet body = %#v, want %#v", got, currentArrayAtIndexBody)
			}
		case iteration9MigrationName:
			foundNinth = true
			if got, want := ordered[index].snippet.Prefix, "array concat"; got != want {
				t.Fatalf("iteration-9 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "concatenate two arrays"; got != want {
				t.Fatalf("iteration-9 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayConcatBody {
				t.Fatalf("iteration-9 snippet body = %#v, want %#v", got, currentArrayConcatBody)
			}
		case iteration10MigrationName:
			foundTenth = true
			if got, want := ordered[index].snippet.Prefix, "array contains"; got != want {
				t.Fatalf("iteration-10 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "check if the array contains an element"; got != want {
				t.Fatalf("iteration-10 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentArrayContainsBody) {
				t.Fatalf("iteration-10 snippet body = %#v, want %#v", got, currentArrayContainsBody)
			}
		case iteration11MigrationName:
			foundEleventh = true
			if got, want := ordered[index].snippet.Prefix, "array declare"; got != want {
				t.Fatalf("iteration-11 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "declare an array"; got != want {
				t.Fatalf("iteration-11 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentArrayDeclareBody) {
				t.Fatalf("iteration-11 snippet body = %#v, want %#v", got, currentArrayDeclareBody)
			}
		case iteration12MigrationName:
			foundTwelfth = true
			if got, want := ordered[index].snippet.Prefix, "array delete at"; got != want {
				t.Fatalf("iteration-12 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "delete element at index from array"; got != want {
				t.Fatalf("iteration-12 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayDeleteAtBody {
				t.Fatalf("iteration-12 snippet body = %#v, want %#v", got, currentArrayDeleteAtBody)
			}
		case iteration13MigrationName:
			foundThirteenth = true
			if got, want := ordered[index].snippet.Prefix, "array delete"; got != want {
				t.Fatalf("iteration-13 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "delete entire array"; got != want {
				t.Fatalf("iteration-13 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayDeleteBody {
				t.Fatalf("iteration-13 snippet body = %#v, want %#v", got, currentArrayDeleteBody)
			}
		case iteration14MigrationName:
			foundFourteenth = true
			if got, want := ordered[index].snippet.Prefix, "array filter"; got != want {
				t.Fatalf("iteration-14 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "filter elements of an array based on given grep pattern"; got != want {
				t.Fatalf("iteration-14 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayFilterBody {
				t.Fatalf("iteration-14 snippet body = %#v, want %#v", got, currentArrayFilterBody)
			}
		case iteration15MigrationName:
			foundFifteenth = true
			wantPrefix := []any{"array iterate", "array forEach"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-15 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "iterate array elements"; got != want {
				t.Fatalf("iteration-15 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentArrayIterateBody) {
				t.Fatalf("iteration-15 snippet body = %#v, want %#v", got, currentArrayIterateBody)
			}
		case iteration16MigrationName:
			foundSixteenth = true
			if got, want := ordered[index].snippet.Prefix, "array length"; got != want {
				t.Fatalf("iteration-16 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "length of an array"; got != want {
				t.Fatalf("iteration-16 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayLengthBody {
				t.Fatalf("iteration-16 snippet body = %#v, want %#v", got, currentArrayLengthBody)
			}
		case iteration17MigrationName:
			foundSeventeenth = true
			wantPrefix := []any{"array print", "echo array"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-17 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "iterate array elements"; got != want {
				t.Fatalf("iteration-17 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayPrintBody {
				t.Fatalf("iteration-17 snippet body = %#v, want %#v", got, currentArrayPrintBody)
			}
		case iteration18MigrationName:
			foundEighteenth = true
			wantPrefix := []any{"array push", "array add"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-18 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "push new item to the end of array"; got != want {
				t.Fatalf("iteration-18 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayPushBody {
				t.Fatalf("iteration-18 snippet body = %#v, want %#v", got, currentArrayPushBody)
			}
		case iteration19MigrationName:
			foundNineteenth = true
			wantPrefix := []any{"array slice", "array range"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-19 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "n elements of an array from specified index (zero based)"; got != want {
				t.Fatalf("iteration-19 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayRangeBody {
				t.Fatalf("iteration-19 snippet body = %#v, want %#v", got, currentArrayRangeBody)
			}
		case iteration20MigrationName:
			foundTwentieth = true
			if got, want := ordered[index].snippet.Prefix, "array replace"; got != want {
				t.Fatalf("iteration-20 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "find and replace elements in array using regex"; got != want {
				t.Fatalf("iteration-20 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArrayReplaceBody {
				t.Fatalf("iteration-20 snippet body = %#v, want %#v", got, currentArrayReplaceBody)
			}
		case iteration21MigrationName:
			foundTwentyFirst = true
			if got, want := ordered[index].snippet.Prefix, "array reverse"; got != want {
				t.Fatalf("iteration-21 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "reverse order of array elements"; got != want {
				t.Fatalf("iteration-21 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentArrayReverseBody) {
				t.Fatalf("iteration-21 snippet body = %#v, want %#v", got, currentArrayReverseBody)
			}
		case iteration22MigrationName:
			foundTwentySecond = true
			if got, want := ordered[index].snippet.Prefix, "array set element"; got != want {
				t.Fatalf("iteration-22 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "set array element at specified index"; got != want {
				t.Fatalf("iteration-22 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentArraySetElementAtBody {
				t.Fatalf("iteration-22 snippet body = %#v, want %#v", got, currentArraySetElementAtBody)
			}
		case iteration23MigrationName:
			foundTwentyThird = true
			wantPrefix := []any{"command failure check", "cmd failure check"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-23 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "check if last command failed"; got != want {
				t.Fatalf("iteration-23 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentCommandFailureCheckBody) {
				t.Fatalf("iteration-23 snippet body = %#v, want %#v", got, currentCommandFailureCheckBody)
			}
		case iteration24MigrationName:
			foundTwentyFourth = true
			wantPrefix := []any{"hide command error", "don't show command error"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-24 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "If a command fails don't show error (suppress stderr)"; got != want {
				t.Fatalf("iteration-24 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCommandHideErrorBody {
				t.Fatalf("iteration-24 snippet body = %#v, want %#v", got, currentCommandHideErrorBody)
			}
		case iteration25MigrationName:
			foundTwentyFifth = true
			wantPrefix := []any{"if command exists", "if cmd exists"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-25 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "check if command exists"; got != want {
				t.Fatalf("iteration-25 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentCommandIfExistsBody) {
				t.Fatalf("iteration-25 snippet body = %#v, want %#v", got, currentCommandIfExistsBody)
			}
		case iteration26MigrationName:
			foundTwentySixth = true
			wantPrefix := []any{"command nice", "cmd nice"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-26 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "run command with desired privilege. n: -20 (highest priority) to 19 (lowest priority)"; got != want {
				t.Fatalf("iteration-26 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCommandNiceBody {
				t.Fatalf("iteration-26 snippet body = %#v, want %#v", got, currentCommandNiceBody)
			}
		case iteration27MigrationName:
			foundTwentySeventh = true
			wantPrefix := []any{"command renice", "cmd renice"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-27 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "change running process priority. n: -20 (highest priority) to 19 (lowest priority)"; got != want {
				t.Fatalf("iteration-27 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCommandReniceBody {
				t.Fatalf("iteration-27 snippet body = %#v, want %#v", got, currentCommandReniceBody)
			}
		case iteration28MigrationName:
			foundTwentyEighth = true
			wantPrefix := []any{"command", "cmd", "command substitution", "cmd substitution"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-28 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "run command (command substitution)"; got != want {
				t.Fatalf("iteration-28 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCommandRunBody {
				t.Fatalf("iteration-28 snippet body = %#v, want %#v", got, currentCommandRunBody)
			}
		case iteration29MigrationName:
			foundTwentyNinth = true
			wantPrefix := []any{"command substitution"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-29 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "run command (command substitution)"; got != want {
				t.Fatalf("iteration-29 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCommandSubstitutionBody {
				t.Fatalf("iteration-29 snippet body = %#v, want %#v", got, currentCommandSubstitutionBody)
			}
		case iteration30MigrationName:
			foundThirtieth = true
			wantPrefix := []any{"command success check", "cmd success check"}
			if !reflect.DeepEqual(ordered[index].snippet.Prefix, wantPrefix) {
				t.Fatalf("iteration-30 snippet prefixes = %#v, want %#v", ordered[index].snippet.Prefix, wantPrefix)
			}
			if got, want := ordered[index].snippet.Description, "check if last command succeed"; got != want {
				t.Fatalf("iteration-30 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentCommandSuccessCheckBody) {
				t.Fatalf("iteration-30 snippet body = %#v, want %#v", got, currentCommandSuccessCheckBody)
			}
		case iteration31MigrationName:
			foundThirtyFirst = true
			if got, want := ordered[index].snippet.Prefix, "crypto base64 decode"; got != want {
				t.Fatalf("iteration-31 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "decode variable from base64"; got != want {
				t.Fatalf("iteration-31 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCryptographyBase64DecodeBody {
				t.Fatalf("iteration-31 snippet body = %#v, want %#v", got, currentCryptographyBase64DecodeBody)
			}
		case iteration32MigrationName:
			foundThirtySecond = true
			if got, want := ordered[index].snippet.Prefix, "crypto base64 encode"; got != want {
				t.Fatalf("iteration-32 snippet prefix = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "encode variable to base64"; got != want {
				t.Fatalf("iteration-32 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCryptographyBase64EncodeBody {
				t.Fatalf("iteration-32 snippet body = %#v, want %#v", got, currentCryptographyBase64EncodeBody)
			}
		case iteration33MigrationName:
			foundThirtyThird = true
			if got, want := ordered[index].snippet.Prefix, "crypto hash"; got != want {
				t.Fatalf("iteration-33 snippet prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "compute hash of variable (md5, sha, sha1, sha224, sha256, sha384, sha512)"; got != want {
				t.Fatalf("iteration-33 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentCryptographyHashBody {
				t.Fatalf("iteration-33 snippet body = %#v, want %#v", got, currentCryptographyHashBody)
			}
		case iteration34MigrationName:
			foundThirtyFourth = true
			if got, want := ordered[index].snippet.Prefix, "date now short"; got != want {
				t.Fatalf("iteration-34 snippet prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "yyyy/mm/dd"; got != want {
				t.Fatalf("iteration-34 snippet description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDateNowShortBody {
				t.Fatalf("iteration-34 snippet body = %#v, want %#v", got, currentDateNowShortBody)
			}
			ordered[index].snippet.Body = v6DateNowShortBody
		}
	}
	if !foundFirst {
		t.Fatalf("approved snippet %q is missing", approvedMigrationName)
	}
	if !foundSecond {
		t.Fatalf("iteration-2 snippet %q is missing", iteration2MigrationName)
	}
	if !foundThird {
		t.Fatalf("iteration-3 snippet %q is missing", iteration3MigrationName)
	}
	if !foundFourth {
		t.Fatalf("iteration-4 snippet %q is missing", iteration4MigrationName)
	}
	if !foundFifth {
		t.Fatalf("iteration-5 snippet %q is missing", iteration5MigrationName)
	}
	if !foundSixth {
		t.Fatalf("iteration-6 snippet %q is missing", iteration6MigrationName)
	}
	if !foundSeventh {
		t.Fatalf("iteration-7 snippet %q is missing", iteration7MigrationName)
	}
	if !foundEighth {
		t.Fatalf("iteration-8 snippet %q is missing", iteration8MigrationName)
	}
	if !foundNinth {
		t.Fatalf("iteration-9 snippet %q is missing", iteration9MigrationName)
	}
	if !foundTenth {
		t.Fatalf("iteration-10 snippet %q is missing", iteration10MigrationName)
	}
	if !foundEleventh {
		t.Fatalf("iteration-11 snippet %q is missing", iteration11MigrationName)
	}
	if !foundTwelfth {
		t.Fatalf("iteration-12 snippet %q is missing", iteration12MigrationName)
	}
	if !foundThirteenth {
		t.Fatalf("iteration-13 snippet %q is missing", iteration13MigrationName)
	}
	if !foundFourteenth {
		t.Fatalf("iteration-14 snippet %q is missing", iteration14MigrationName)
	}
	if !foundFifteenth {
		t.Fatalf("iteration-15 snippet %q is missing", iteration15MigrationName)
	}
	if !foundSixteenth {
		t.Fatalf("iteration-16 snippet %q is missing", iteration16MigrationName)
	}
	if !foundSeventeenth {
		t.Fatalf("iteration-17 snippet %q is missing", iteration17MigrationName)
	}
	if !foundEighteenth {
		t.Fatalf("iteration-18 snippet %q is missing", iteration18MigrationName)
	}
	if !foundNineteenth {
		t.Fatalf("iteration-19 snippet %q is missing", iteration19MigrationName)
	}
	if !foundTwentieth {
		t.Fatalf("iteration-20 snippet %q is missing", iteration20MigrationName)
	}
	if !foundTwentyFirst {
		t.Fatalf("iteration-21 snippet %q is missing", iteration21MigrationName)
	}
	if !foundTwentySecond {
		t.Fatalf("iteration-22 snippet %q is missing", iteration22MigrationName)
	}
	if !foundTwentyThird {
		t.Fatalf("iteration-23 snippet %q is missing", iteration23MigrationName)
	}
	if !foundTwentyFourth {
		t.Fatalf("iteration-24 snippet %q is missing", iteration24MigrationName)
	}
	if !foundTwentyFifth {
		t.Fatalf("iteration-25 snippet %q is missing", iteration25MigrationName)
	}
	if !foundTwentySixth {
		t.Fatalf("iteration-26 snippet %q is missing", iteration26MigrationName)
	}
	if !foundTwentySeventh {
		t.Fatalf("iteration-27 snippet %q is missing", iteration27MigrationName)
	}
	if !foundTwentyEighth {
		t.Fatalf("iteration-28 snippet %q is missing", iteration28MigrationName)
	}
	if !foundTwentyNinth {
		t.Fatalf("iteration-29 snippet %q is missing", iteration29MigrationName)
	}
	if !foundThirtieth {
		t.Fatalf("iteration-30 snippet %q is missing", iteration30MigrationName)
	}
	if !foundThirtyFirst {
		t.Fatalf("iteration-31 snippet %q is missing", iteration31MigrationName)
	}
	if !foundThirtySecond {
		t.Fatalf("iteration-32 snippet %q is missing", iteration32MigrationName)
	}
	if !foundThirtyThird {
		t.Fatalf("iteration-33 snippet %q is missing", iteration33MigrationName)
	}
	if !foundThirtyFourth {
		t.Fatalf("iteration-34 snippet %q is missing", iteration34MigrationName)
	}

	reconstructedIteration33, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-33 snippets", reconstructedIteration33, currentV7Iteration33SHA256)
	assertSHA256(t, "approved v7 iteration-33 commands", renderDocumentation(ordered), currentV7Iteration33Commands)

	for index := range ordered {
		if ordered[index].name == iteration33MigrationName {
			ordered[index].snippet.Body = v6CryptographyHashBody
		}
	}

	reconstructedIteration32, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-32 snippets", reconstructedIteration32, currentV7Iteration32SHA256)
	assertSHA256(t, "approved v7 iteration-32 commands", renderDocumentation(ordered), currentV7Iteration32Commands)

	for index := range ordered {
		if ordered[index].name == iteration32MigrationName {
			ordered[index].snippet.Body = v6CryptographyBase64EncodeBody
		}
	}

	reconstructedIteration31, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-31 snippets", reconstructedIteration31, currentV7Iteration31SHA256)
	assertSHA256(t, "approved v7 iteration-31 commands", renderDocumentation(ordered), currentV7Iteration31Commands)

	for index := range ordered {
		if ordered[index].name == iteration31MigrationName {
			ordered[index].snippet.Body = v6CryptographyBase64DecodeBody
		}
	}

	reconstructedIteration30, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-30 snippets", reconstructedIteration30, currentV7Iteration30SHA256)
	assertSHA256(t, "approved v7 iteration-30 commands", renderDocumentation(ordered), currentV7Iteration30Commands)

	for index := range ordered {
		if ordered[index].name == iteration30MigrationName {
			ordered[index].snippet.Body = v6CommandSuccessCheckBody
		}
	}

	reconstructedIteration29, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-29 snippets", reconstructedIteration29, currentV7Iteration29SHA256)
	assertSHA256(t, "approved v7 iteration-29 commands", renderDocumentation(ordered), currentV7Iteration29Commands)

	for index := range ordered {
		if ordered[index].name == iteration29MigrationName {
			ordered[index].snippet.Body = v6CommandSubstitutionBody
		}
	}

	reconstructedIteration28, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-28 snippets", reconstructedIteration28, currentV7Iteration28SHA256)
	assertSHA256(t, "approved v7 iteration-28 commands", renderDocumentation(ordered), currentV7Iteration28Commands)

	for index := range ordered {
		if ordered[index].name == iteration28MigrationName {
			ordered[index].snippet.Body = v6CommandRunBody
		}
	}

	reconstructedIteration27, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-27 snippets", reconstructedIteration27, currentV7Iteration27SHA256)
	assertSHA256(t, "approved v7 iteration-27 commands", renderDocumentation(ordered), currentV7Iteration27Commands)

	for index := range ordered {
		if ordered[index].name == iteration27MigrationName {
			ordered[index].snippet.Body = v6CommandReniceBody
		}
	}

	reconstructedIteration26, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-26 snippets", reconstructedIteration26, currentV7Iteration26SHA256)
	assertSHA256(t, "approved v7 iteration-26 commands", renderDocumentation(ordered), currentV7Iteration26Commands)

	for index := range ordered {
		if ordered[index].name == iteration26MigrationName {
			ordered[index].snippet.Body = v6CommandNiceBody
		}
	}

	reconstructedIteration25, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-25 snippets", reconstructedIteration25, currentV7Iteration25SHA256)
	assertSHA256(t, "approved v7 iteration-25 commands", renderDocumentation(ordered), currentV7Iteration25Commands)

	for index := range ordered {
		if ordered[index].name == iteration25MigrationName {
			ordered[index].snippet.Body = v6CommandIfExistsBody
		}
	}

	reconstructedIteration24, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-24 snippets", reconstructedIteration24, currentV7Iteration24SHA256)
	assertSHA256(t, "approved v7 iteration-24 commands", renderDocumentation(ordered), currentV7Iteration24Commands)

	for index := range ordered {
		if ordered[index].name == iteration24MigrationName {
			ordered[index].snippet.Body = v6CommandHideErrorBody
		}
	}

	reconstructedIteration23, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-23 snippets", reconstructedIteration23, currentV7Iteration23SHA256)
	assertSHA256(t, "approved v7 iteration-23 commands", renderDocumentation(ordered), currentV7Iteration23Commands)

	for index := range ordered {
		if ordered[index].name == iteration23MigrationName {
			ordered[index].snippet.Body = v6CommandFailureCheckBody
		}
	}

	reconstructedIteration22, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-22 snippets", reconstructedIteration22, currentV7Iteration22SHA256)
	assertSHA256(t, "approved v7 iteration-22 commands", renderDocumentation(ordered), currentV7Iteration22Commands)

	for index := range ordered {
		if ordered[index].name == iteration22MigrationName {
			ordered[index].snippet.Body = v6ArraySetElementAtBody
		}
	}

	reconstructedIteration21, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-21 snippets", reconstructedIteration21, currentV7Iteration21SHA256)
	assertSHA256(t, "approved v7 iteration-21 commands", renderDocumentation(ordered), currentV7Iteration21Commands)

	for index := range ordered {
		if ordered[index].name == iteration21MigrationName {
			ordered[index].snippet.Body = v6ArrayReverseBody
		}
	}

	reconstructedIteration20, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-20 snippets", reconstructedIteration20, approvedV7Iteration20SHA256)
	assertSHA256(t, "approved v7 iteration-20 commands", renderDocumentation(ordered), approvedV7Iteration20Commands)

	for index := range ordered {
		if ordered[index].name == iteration20MigrationName {
			ordered[index].snippet.Body = v6ArrayReplaceBody
		}
	}

	reconstructedIteration19, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-19 snippets", reconstructedIteration19, approvedV7Iteration19SHA256)
	assertSHA256(t, "approved v7 iteration-19 commands", renderDocumentation(ordered), approvedV7Iteration19Commands)

	for index := range ordered {
		if ordered[index].name == iteration19MigrationName {
			ordered[index].snippet.Body = v6ArrayRangeBody
		}
	}

	reconstructedIteration18, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-18 snippets", reconstructedIteration18, approvedV7Iteration18SHA256)
	assertSHA256(t, "approved v7 iteration-18 commands", renderDocumentation(ordered), approvedV7Iteration18Commands)

	for index := range ordered {
		if ordered[index].name == iteration18MigrationName {
			ordered[index].snippet.Body = v6ArrayPushBody
		}
	}

	reconstructedIteration17, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-17 snippets", reconstructedIteration17, approvedV7Iteration17SHA256)
	assertSHA256(t, "approved v7 iteration-17 commands", renderDocumentation(ordered), approvedV7Iteration17Commands)

	for index := range ordered {
		if ordered[index].name == iteration17MigrationName {
			ordered[index].snippet.Body = v6ArrayPrintBody
		}
	}

	reconstructedIteration16, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-16 snippets", reconstructedIteration16, approvedV7Iteration16SHA256)
	assertSHA256(t, "approved v7 iteration-16 commands", renderDocumentation(ordered), approvedV7Iteration16Commands)

	for index := range ordered {
		if ordered[index].name == iteration16MigrationName {
			ordered[index].snippet.Body = v6ArrayLengthBody
		}
	}

	reconstructedIteration15, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-15 snippets", reconstructedIteration15, approvedV7Iteration15SHA256)
	assertSHA256(t, "approved v7 iteration-15 commands", renderDocumentation(ordered), approvedV7Iteration15Commands)

	for index := range ordered {
		if ordered[index].name == iteration15MigrationName {
			ordered[index].snippet.Body = v6ArrayIterateBody
		}
	}

	reconstructedIteration14, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-14 snippets", reconstructedIteration14, approvedV7Iteration14SHA256)
	assertSHA256(t, "approved v7 iteration-14 commands", renderDocumentation(ordered), approvedV7Iteration14Commands)

	for index := range ordered {
		if ordered[index].name == iteration14MigrationName {
			ordered[index].snippet.Body = v6ArrayFilterBody
		}
	}

	reconstructedIteration13, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-13 snippets", reconstructedIteration13, approvedV7Iteration13SHA256)
	assertSHA256(t, "approved v7 iteration-13 commands", renderDocumentation(ordered), approvedV7Iteration13Commands)

	for index := range ordered {
		if ordered[index].name == iteration13MigrationName {
			ordered[index].snippet.Body = v6ArrayDeleteBody
		}
	}

	reconstructedIteration12, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-12 snippets", reconstructedIteration12, approvedV7Iteration12SHA256)
	assertSHA256(t, "approved v7 iteration-12 commands", renderDocumentation(ordered), approvedV7Iteration12Commands)

	for index := range ordered {
		if ordered[index].name == iteration12MigrationName {
			ordered[index].snippet.Body = v6ArrayDeleteAtBody
		}
	}

	reconstructedIteration11, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-11 snippets", reconstructedIteration11, approvedV7Iteration11SHA256)
	assertSHA256(t, "approved v7 iteration-11 commands", renderDocumentation(ordered), approvedV7Iteration11Commands)

	for index := range ordered {
		if ordered[index].name == iteration11MigrationName {
			ordered[index].snippet.Body = v6ArrayDeclareBody
		}
	}

	reconstructedIteration10, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-10 snippets", reconstructedIteration10, approvedV7Iteration10SHA256)
	assertSHA256(t, "approved v7 iteration-10 commands", renderDocumentation(ordered), approvedV7Iteration10Commands)

	for index := range ordered {
		if ordered[index].name == iteration10MigrationName {
			ordered[index].snippet.Body = v6ArrayContainsBody
		}
	}

	reconstructedIteration9, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-9 snippets", reconstructedIteration9, approvedV7Iteration9SHA256)
	assertSHA256(t, "approved v7 iteration-9 commands", renderDocumentation(ordered), approvedV7Iteration9Commands)

	for index := range ordered {
		if ordered[index].name == iteration9MigrationName {
			ordered[index].snippet.Body = v6ArrayConcatBody
		}
	}

	reconstructedIteration8, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-8 snippets", reconstructedIteration8, approvedV7Iteration8SHA256)
	assertSHA256(t, "approved v7 iteration-8 commands", renderDocumentation(ordered), approvedV7Iteration8Commands)

	for index := range ordered {
		if ordered[index].name == iteration8MigrationName {
			ordered[index].snippet.Body = v6ArrayAtIndexBody
		}
	}

	reconstructedIteration7, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-7 snippets", reconstructedIteration7, approvedV7Iteration7SHA256)
	assertSHA256(t, "approved v7 iteration-7 commands", renderDocumentation(ordered), approvedV7Iteration7Commands)

	for index := range ordered {
		if ordered[index].name == iteration7MigrationName {
			ordered[index].snippet.Body = v6ArrayAllElementsBody
		}
	}

	reconstructedIteration6, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-6 snippets", reconstructedIteration6, approvedV7Iteration6SHA256)
	assertSHA256(t, "approved v7 iteration-6 commands", renderDocumentation(ordered), approvedV7Iteration6Commands)

	for index := range ordered {
		if ordered[index].name == iteration6MigrationName {
			ordered[index].snippet.Body = v6DecompressUnzipBody
		}
	}

	reconstructedIteration5, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-5 snippets", reconstructedIteration5, approvedV7Iteration5SHA256)
	assertSHA256(t, "approved v7 iteration-5 commands", renderDocumentation(ordered), approvedV7Iteration5Commands)

	for index := range ordered {
		if ordered[index].name == iteration5MigrationName {
			ordered[index].snippet.Body = v6DecompressTarXzBody
		}
	}

	reconstructedIteration4, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-4 snippets", reconstructedIteration4, approvedV7Iteration4SHA256)
	assertSHA256(t, "approved v7 iteration-4 commands", renderDocumentation(ordered), approvedV7Iteration4Commands)

	for index := range ordered {
		if ordered[index].name == iteration4MigrationName {
			ordered[index].snippet.Body = v6DecompressTarGzBody
		}
	}

	reconstructedIteration3, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-3 snippets", reconstructedIteration3, approvedV7Iteration3SHA256)
	assertSHA256(t, "approved v7 iteration-3 commands", renderDocumentation(ordered), approvedV7Iteration3Commands)

	for index := range ordered {
		if ordered[index].name == iteration3MigrationName {
			ordered[index].snippet.Body = v6CompressZipBody
		}
	}

	reconstructedIteration2, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-2 snippets", reconstructedIteration2, approvedV7Iteration2SHA256)
	assertSHA256(t, "approved v7 iteration-2 commands", renderDocumentation(ordered), approvedV7Iteration2Commands)

	for index := range ordered {
		if ordered[index].name == iteration2MigrationName {
			ordered[index].snippet.Body = v6CompressTarXzBody
		}
	}

	reconstructedIteration1, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-1 snippets", reconstructedIteration1, approvedV7Iteration1SHA256)
	assertSHA256(t, "approved v7 iteration-1 commands", renderDocumentation(ordered), approvedV7Iteration1Commands)

	for index := range ordered {
		if ordered[index].name == approvedMigrationName {
			ordered[index].snippet.Body = v6CompressTarGzBody
		}
	}
	reconstructedV6, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed historical v6 snippets", reconstructedV6, historicalV6SnippetSHA256)
	assertSHA256(t, "historical v6 commands", renderDocumentation(ordered), historicalV6CommandsSHA256)

	changed := 0
	for index := range ordered {
		if !reflect.DeepEqual(currentOrdered[index].snippet, ordered[index].snippet) {
			changed++
		}
	}
	if changed != 34 {
		t.Fatalf("snippets differing from historical v6 = %d, want 34 (244 unchanged)", changed)
	}
}

func TestCompressTarGzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentCompressTarGzBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `tar -cf "${archive_path}.tar" "${source_path}"`) {
		t.Fatal("snippet must use the portable tar -c and -f options")
	}
	for _, nonPortable := range []string{" -z", " --", " -C"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
}

func TestCompressTarGzBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "tar")
	requireCommand(t, "gzip")

	script := runnableGeneratedCompressTarGz(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters and integrity", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source dir [brackets] #dollar$"
		source := filepath.Join(root, sourceName)
		mustMkdirAll(t, filepath.Join(source, "empty dir"))
		mustWriteFile(t, filepath.Join(source, "file name [1] #$.txt"), []byte("shellman archive integrity\n"))
		archiveBase := filepath.Join(root, "output dir", "archive [1] #$")
		mustMkdirAll(t, filepath.Dir(archiveBase))

		run(t, root, "sh", scriptPath(t, script), archiveBase, sourceName)
		archive := archiveBase + ".tar.gz"
		run(t, root, "gzip", "-t", archive)
		listing := run(t, root, "tar", "-tzf", archive)
		for _, entry := range []string{sourceName + "/", sourceName + "/empty dir/", sourceName + "/file name [1] #$.txt"} {
			if !strings.Contains(listing, entry+"\n") {
				t.Fatalf("archive listing lacks %q:\n%s", entry, listing)
			}
		}

		extract := filepath.Join(root, "extracted")
		mustMkdirAll(t, extract)
		run(t, extract, "tar", "-xzf", archive)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "shellman archive integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "empty source"))
		archiveBase := filepath.Join(root, "empty archive")
		run(t, root, "sh", scriptPath(t, script), archiveBase, "empty source")
		listing := run(t, root, "tar", "-tzf", archiveBase+".tar.gz")
		if listing != "empty source/\n" {
			t.Fatalf("empty-directory archive listing = %q", listing)
		}
	})

	t.Run("missing source fails without gzip output", func(t *testing.T) {
		root := t.TempDir()
		archiveBase := filepath.Join(root, "missing archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "does not exist")
		if err == nil {
			t.Fatal("missing source unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".tar.gz"); !os.IsNotExist(err) {
			t.Fatalf("gzip archive exists after tar failure: %v", err)
		}
	})

	t.Run("existing destination is overwritten", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("replacement\n"))
		archiveBase := filepath.Join(root, "existing")
		mustWriteFile(t, archiveBase+".tar.gz", []byte("old invalid archive"))
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		run(t, root, "gzip", "-t", archiveBase+".tar.gz")
		listing := run(t, root, "tar", "-tzf", archiveBase+".tar.gz")
		if listing != "source.txt\n" {
			t.Fatalf("replacement archive listing = %q", listing)
		}
	})

	t.Run("tar failure prevents gzip", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 19\n")
		marker := filepath.Join(root, "gzip-ran")
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\n: > \"$SHELLMAN_GZIP_MARKER\"\nexit 0\n")
		environment := append(os.Environ(), "PATH="+fakeBin, "SHELLMAN_GZIP_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), "source")
		if exitCode(err) != 19 {
			t.Fatalf("exit status = %d, want tar status 19: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("gzip ran after tar failure: %v", err)
		}
	})

	t.Run("gzip failure leaves intermediate tar", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("content\n"))
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\nexit 23\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		archiveBase := filepath.Join(root, "gzip failure")
		_, err := runCommand(root, environment, "sh", scriptFile, archiveBase, "source.txt")
		if exitCode(err) != 23 {
			t.Fatalf("exit status = %d, want gzip status 23: %v", exitCode(err), err)
		}
		if _, err := os.Stat(archiveBase + ".tar"); err != nil {
			t.Fatalf("tar intermediate not retained after gzip failure: %v", err)
		}
		if _, err := os.Stat(archiveBase + ".tar.gz"); !os.IsNotExist(err) {
			t.Fatalf("gzip output exists after simulated failure: %v", err)
		}
	})

	t.Run("source beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-source"))
		mustWriteFile(t, filepath.Join(root, "-source", "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "hyphen archive")
		run(t, root, "sh", scriptFile, archiveBase, "-source")
		listing := run(t, root, "tar", "-tzf", archiveBase+".tar.gz")
		if !strings.Contains(listing, "./-source/file.txt\n") {
			t.Fatalf("hyphen-source archive listing = %q", listing)
		}
	})

	t.Run("destination inside source directory", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source"
		mustMkdirAll(t, filepath.Join(root, sourceName))
		mustWriteFile(t, filepath.Join(root, sourceName, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, sourceName, "inside")
		run(t, root, "sh", scriptFile, archiveBase, sourceName)
		archive := archiveBase + ".tar.gz"
		run(t, root, "gzip", "-t", archive)
		listing := run(t, root, "tar", "-tzf", archive)
		if !strings.Contains(listing, "source/file.txt\n") {
			t.Fatalf("inside-source archive lacks input file: %q", listing)
		}
		if strings.Contains(listing, "source/inside.tar") || strings.Contains(listing, "source/inside.tar.gz") {
			t.Fatalf("archive included itself: %q", listing)
		}
	})
}

func runnableGeneratedCompressTarGz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[approvedMigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", approvedMigrationName)
	}
	if body != currentCompressTarGzBody {
		t.Fatal("generated snippet body differs from the approved body")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`, `$2`)
}

func TestCompressTarXzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentCompressTarXzBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `tar -cf "${archive_path}.tar" "${source_path}"`) {
		t.Fatal("snippet must use the portable tar -c and -f options")
	}
	for _, nonPortable := range []string{" -J", " --", " -C"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
}

func TestCompressTarXzBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "tar")
	requireCommand(t, "xz")

	script := runnableGeneratedCompressTarXz(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters and integrity", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source dir [brackets] #dollar$"
		source := filepath.Join(root, sourceName)
		mustMkdirAll(t, filepath.Join(source, "empty dir"))
		mustWriteFile(t, filepath.Join(source, "file name [1] #$.txt"), []byte("shellman xz archive integrity\n"))
		archiveBase := filepath.Join(root, "output dir", "archive [1] #$")
		mustMkdirAll(t, filepath.Dir(archiveBase))

		run(t, root, "sh", scriptPath(t, script), archiveBase, sourceName)
		archive := archiveBase + ".tar.xz"
		run(t, root, "xz", "-t", archive)
		listing := run(t, root, "tar", "-tJf", archive)
		for _, entry := range []string{sourceName + "/", sourceName + "/empty dir/", sourceName + "/file name [1] #$.txt"} {
			if !strings.Contains(listing, entry+"\n") {
				t.Fatalf("archive listing lacks %q:\n%s", entry, listing)
			}
		}

		extract := filepath.Join(root, "extracted")
		mustMkdirAll(t, extract)
		run(t, extract, "tar", "-xJf", archive)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "shellman xz archive integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "empty source"))
		archiveBase := filepath.Join(root, "empty archive")
		run(t, root, "sh", scriptPath(t, script), archiveBase, "empty source")
		listing := run(t, root, "tar", "-tJf", archiveBase+".tar.xz")
		if listing != "empty source/\n" {
			t.Fatalf("empty-directory archive listing = %q", listing)
		}
	})

	t.Run("missing source fails without xz output", func(t *testing.T) {
		root := t.TempDir()
		archiveBase := filepath.Join(root, "missing archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "does not exist")
		if err == nil {
			t.Fatal("missing source unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".tar.xz"); !os.IsNotExist(err) {
			t.Fatalf("xz archive exists after tar failure: %v", err)
		}
	})

	t.Run("existing destination is overwritten", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("replacement\n"))
		archiveBase := filepath.Join(root, "existing")
		mustWriteFile(t, archiveBase+".tar.xz", []byte("old invalid archive"))
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		run(t, root, "xz", "-t", archiveBase+".tar.xz")
		listing := run(t, root, "tar", "-tJf", archiveBase+".tar.xz")
		if listing != "source.txt\n" {
			t.Fatalf("replacement archive listing = %q", listing)
		}
	})

	t.Run("tar failure prevents xz", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 29\n")
		marker := filepath.Join(root, "xz-ran")
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\n: > \"$SHELLMAN_XZ_MARKER\"\nexit 0\n")
		environment := append(os.Environ(), "PATH="+fakeBin, "SHELLMAN_XZ_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), "source")
		if exitCode(err) != 29 {
			t.Fatalf("exit status = %d, want tar status 29: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("xz ran after tar failure: %v", err)
		}
	})

	t.Run("xz failure leaves intermediate tar", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("content\n"))
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nexit 31\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		archiveBase := filepath.Join(root, "xz failure")
		_, err := runCommand(root, environment, "sh", scriptFile, archiveBase, "source.txt")
		if exitCode(err) != 31 {
			t.Fatalf("exit status = %d, want xz status 31: %v", exitCode(err), err)
		}
		if _, err := os.Stat(archiveBase + ".tar"); err != nil {
			t.Fatalf("tar intermediate not retained after xz failure: %v", err)
		}
		if _, err := os.Stat(archiveBase + ".tar.xz"); !os.IsNotExist(err) {
			t.Fatalf("xz output exists after simulated failure: %v", err)
		}
	})

	t.Run("source beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-source"))
		mustWriteFile(t, filepath.Join(root, "-source", "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "hyphen archive")
		run(t, root, "sh", scriptFile, archiveBase, "-source")
		listing := run(t, root, "tar", "-tJf", archiveBase+".tar.xz")
		if !strings.Contains(listing, "./-source/file.txt\n") {
			t.Fatalf("hyphen-source archive listing = %q", listing)
		}
	})

	t.Run("destination inside source directory", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source"
		mustMkdirAll(t, filepath.Join(root, sourceName))
		mustWriteFile(t, filepath.Join(root, sourceName, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, sourceName, "inside")
		run(t, root, "sh", scriptFile, archiveBase, sourceName)
		archive := archiveBase + ".tar.xz"
		run(t, root, "xz", "-t", archive)
		listing := run(t, root, "tar", "-tJf", archive)
		if !strings.Contains(listing, "source/file.txt\n") {
			t.Fatalf("inside-source archive lacks input file: %q", listing)
		}
		if strings.Contains(listing, "source/inside.tar") || strings.Contains(listing, "source/inside.tar.xz") {
			t.Fatalf("archive included itself: %q", listing)
		}
	})
}

func runnableGeneratedCompressTarXz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration2MigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", iteration2MigrationName)
	}
	if body != currentCompressTarXzBody {
		t.Fatal("generated iteration-2 snippet body differs from the approved candidate body")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/directory-or-file, "${pathToDirectoryOrFile}"|}`, `$2`)
}

func TestCompressZipPlaceholderContractAndOptionSafety(t *testing.T) {
	body := currentCompressZipBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/path/to/directory-or-file,"${pathToDirectoryOrFile}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `zip -rq "${archive_path}.zip" "${source_path}"`) {
		t.Fatal("snippet must quote the archive and source paths")
	}
	if strings.Contains(body, " --") {
		t.Fatal("snippet must not rely on non-portable end-of-options syntax")
	}
}

func TestCompressZipBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "zip")
	requireCommand(t, "unzip")

	script := runnableGeneratedCompressZip(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters and integrity", func(t *testing.T) {
		root := t.TempDir()
		sourceName := "source dir [brackets] #dollar$"
		source := filepath.Join(root, sourceName)
		mustMkdirAll(t, filepath.Join(source, "empty dir"))
		mustWriteFile(t, filepath.Join(source, "file name [1] #$.txt"), []byte("shellman zip archive integrity\n"))
		archiveBase := filepath.Join(root, "output dir", "archive [1] #$")
		mustMkdirAll(t, filepath.Dir(archiveBase))

		run(t, root, "sh", scriptPath(t, script), archiveBase, sourceName)
		archive := archiveBase + ".zip"
		run(t, root, "unzip", "-t", archive)
		listing := run(t, root, "unzip", "-Z1", archive)
		for _, entry := range []string{sourceName + "/", sourceName + "/empty dir/", sourceName + "/file name [1] #$.txt"} {
			if !strings.Contains(listing, entry+"\n") {
				t.Fatalf("archive listing lacks %q:\n%s", entry, listing)
			}
		}

		extract := filepath.Join(root, "extracted")
		mustMkdirAll(t, extract)
		run(t, root, "unzip", "-q", archive, "-d", extract)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "shellman zip archive integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "empty source"))
		archiveBase := filepath.Join(root, "empty archive")
		run(t, root, "sh", scriptFile, archiveBase, "empty source")
		listing := run(t, root, "unzip", "-Z1", archiveBase+".zip")
		if listing != "empty source/\n" {
			t.Fatalf("empty-directory archive listing = %q", listing)
		}
	})

	t.Run("missing source fails without output", func(t *testing.T) {
		root := t.TempDir()
		archiveBase := filepath.Join(root, "missing archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "does not exist")
		if err == nil {
			t.Fatal("missing source unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".zip"); !os.IsNotExist(err) {
			t.Fatalf("archive exists after failure: %v", err)
		}
	})

	t.Run("existing archive is updated", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("first\n"))
		archiveBase := filepath.Join(root, "existing")
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("replacement content\n"))
		run(t, root, "sh", scriptFile, archiveBase, "source.txt")
		got := run(t, root, "unzip", "-p", archiveBase+".zip", "source.txt")
		if got != "replacement content\n" {
			t.Fatalf("updated content = %q", got)
		}
	})

	t.Run("zip failure status propagates", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "zip"), "#!/bin/sh\nexit 37\n")
		environment := append(os.Environ(), "PATH="+fakeBin)
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), "source")
		if exitCode(err) != 37 {
			t.Fatalf("exit status = %d, want zip status 37: %v", exitCode(err), err)
		}
	})

	t.Run("archive and source beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-source"))
		mustWriteFile(t, filepath.Join(root, "-source", "file.txt"), []byte("content\n"))
		run(t, root, "sh", scriptFile, "-archive", "-source")
		listing := run(t, root, "unzip", "-Z1", filepath.Join(root, "-archive.zip"))
		if !strings.Contains(listing, "-source/file.txt\n") {
			t.Fatalf("hyphen-path archive listing = %q", listing)
		}
	})

	t.Run("missing destination directory fails", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "source.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "missing", "archive")
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, "source.txt")
		if err == nil {
			t.Fatal("missing destination directory unexpectedly succeeded")
		}
		if _, err := os.Stat(archiveBase + ".zip"); !os.IsNotExist(err) {
			t.Fatalf("archive exists after destination failure: %v", err)
		}
	})
}

func runnableGeneratedCompressZip(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration3MigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", iteration3MigrationName)
	}
	if body != currentCompressZipBody {
		t.Fatal("generated iteration-3 snippet body differs from the approved candidate body")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/directory-or-file,"${pathToDirectoryOrFile}"|}`, `$2`)
}

func TestDecompressTarGzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentDecompressTarGzBody
	for _, placeholder := range []string{
		`${1|/extract/to/path, "${extractToPath}"|}`,
		`${2|/path/to/archive, "${pathToArchive}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `tar -xf "${temporary_directory}/archive.tar"`) {
		t.Fatal("snippet must use portable tar -x and -f options")
	}
	for _, nonPortable := range []string{" -z", " -C", " --"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
	for _, trap := range []string{
		`trap 'status=$?; trap - 0; rm -rf "${temporary_directory}"; exit "${status}"' 0`,
		`trap 'exit 129' HUP`,
		`trap 'exit 130' INT`,
		`trap 'exit 143' TERM`,
	} {
		if !strings.Contains(body, trap) {
			t.Fatalf("snippet lacks separate cleanup/signal trap %q", trap)
		}
	}
}

func TestDecompressTarGzBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "tar")
	requireCommand(t, "gzip")
	requireCommand(t, "mktemp")

	script := runnableGeneratedDecompressTarGz(t)
	scriptFile := scriptPath(t, script)
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("spaces special characters integrity and cleanup", func(t *testing.T) {
		root := t.TempDir()
		staging := filepath.Join(root, "staging")
		sourceName := "source dir [brackets] #dollar$"
		mustMkdirAll(t, filepath.Join(staging, sourceName, "empty dir"))
		mustWriteFile(t, filepath.Join(staging, sourceName, "file name [1] #$.txt"), []byte("decompress integrity\n"))
		archiveBase := filepath.Join(root, "archive [1] #$")
		createTarGz(t, staging, archiveBase, sourceName)
		extract := filepath.Join(root, "extract dir [1] #$")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "temporary root")
		mustMkdirAll(t, temporaryRoot)

		runWithEnvironment(t, root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, archiveBase)
		got, err := os.ReadFile(filepath.Join(extract, sourceName, "file name [1] #$.txt"))
		if err != nil {
			t.Fatal(err)
		}
		if want := "decompress integrity\n"; string(got) != want {
			t.Fatalf("extracted content = %q, want %q", got, want)
		}
		if info, err := os.Stat(filepath.Join(extract, sourceName, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("archive beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createTarGz(t, root, filepath.Join(root, "-archive"), "file.txt")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		run(t, root, "sh", scriptFile, extract, "-archive")
		got, err := os.ReadFile(filepath.Join(extract, "file.txt"))
		if err != nil || string(got) != "content\n" {
			t.Fatalf("hyphen archive extraction = %q, %v", got, err)
		}
	})

	t.Run("missing archive propagates gzip failure and cleans up", func(t *testing.T) {
		root := t.TempDir()
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		_, err := runCommand(root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, filepath.Join(root, "missing"))
		if err == nil {
			t.Fatal("missing archive unexpectedly succeeded")
		}
		assertDirectoryEmpty(t, extract)
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("gzip failure prevents tar and preserves status", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\nexit 41\n")
		marker := filepath.Join(root, "tar-ran")
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\n: > \"$SHELLMAN_TAR_MARKER\"\nexit 0\n")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_TAR_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 41 {
			t.Fatalf("exit status = %d, want gzip status 41: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("tar ran after gzip failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	for _, test := range []struct {
		name       string
		signalFlag string
		exitStatus int
	}{
		{name: "INT cleans up and exits 130", signalFlag: "-INT", exitStatus: 130},
		{name: "TERM cleans up and exits 143", signalFlag: "-TERM", exitStatus: 143},
	} {
		t.Run(test.name, func(t *testing.T) {
			root := t.TempDir()
			fakeBin := filepath.Join(root, "bin")
			mustMkdirAll(t, fakeBin)
			mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\nkill "+test.signalFlag+" \"$PPID\"\nexit 0\n")
			marker := filepath.Join(root, "tar-ran")
			mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\n: > \"$SHELLMAN_TAR_MARKER\"\nexit 0\n")
			extract := filepath.Join(root, "extract")
			mustMkdirAll(t, extract)
			temporaryRoot := filepath.Join(root, "tmp")
			mustMkdirAll(t, temporaryRoot)
			environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_TAR_MARKER="+marker)
			_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
			if exitCode(err) != test.exitStatus {
				t.Fatalf("exit status = %d, want %d: %v", exitCode(err), test.exitStatus, err)
			}
			if _, err := os.Stat(marker); !os.IsNotExist(err) {
				t.Fatalf("tar ran after %s: %v", test.signalFlag, err)
			}
			assertDirectoryEmpty(t, temporaryRoot)
		})
	}

	t.Run("tar failure propagates and cleans up", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createTarGz(t, root, filepath.Join(root, "archive"), "file.txt")
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 43\n")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 43 {
			t.Fatalf("exit status = %d, want tar status 43: %v", exitCode(err), err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("mktemp failure prevents decompression", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "mktemp"), "#!/bin/sh\nexit 47\n")
		marker := filepath.Join(root, "gzip-ran")
		mustExecutable(t, filepath.Join(fakeBin, "gzip"), "#!/bin/sh\n: > \"$SHELLMAN_GZIP_MARKER\"\nexit 0\n")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_GZIP_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 47 {
			t.Fatalf("exit status = %d, want mktemp status 47: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("gzip ran after mktemp failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("missing extraction directory fails after decompression", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "archive")
		createTarGz(t, root, archiveBase, "file.txt")
		temporaryRoot := filepath.Join(root, "tmp")
		mustMkdirAll(t, temporaryRoot)
		_, err := runCommand(root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, filepath.Join(root, "missing"), archiveBase)
		if err == nil {
			t.Fatal("missing extraction directory unexpectedly succeeded")
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})
}

func runnableGeneratedDecompressTarGz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration4MigrationName].Body.(string)
	if !ok {
		t.Fatalf("generated %s body is not a string", iteration4MigrationName)
	}
	if body != currentDecompressTarGzBody {
		t.Fatal("generated iteration-4 snippet body differs from the approved candidate body")
	}
	body = strings.ReplaceAll(body, `${1|/extract/to/path, "${extractToPath}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/archive, "${pathToArchive}"|}`, `$2`)
}

func TestDecompressTarXzPlaceholderContractAndPortableTarOptions(t *testing.T) {
	body := currentDecompressTarXzBody
	for _, placeholder := range []string{
		`${1|/extract/to/path, "${extractToPath}"|}`,
		`${2|/path/to/archive, "${pathToArchive}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	for _, required := range []string{`xz -dc "${archive_path}.tar.xz"`, `tar -xf "${temporary_directory}/archive.tar"`, `trap 'exit 130' INT`, `trap 'exit 143' TERM`} {
		if !strings.Contains(body, required) {
			t.Fatalf("snippet lacks %q", required)
		}
	}
	for _, nonPortable := range []string{" -J", " -C", " --"} {
		if strings.Contains(body, nonPortable) {
			t.Fatalf("snippet uses non-portable tar syntax %q", nonPortable)
		}
	}
}

func TestDecompressTarXzBehavior(t *testing.T) {
	for _, command := range []string{"sh", "tar", "xz", "mktemp"} {
		requireCommand(t, command)
	}
	scriptFile := scriptPath(t, runnableGeneratedDecompressTarXz(t))
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("integrity special paths and normal cleanup", func(t *testing.T) {
		root := t.TempDir()
		staging := filepath.Join(root, "staging")
		source := "source dir [5] #$"
		mustMkdirAll(t, filepath.Join(staging, source, "empty dir"))
		mustWriteFile(t, filepath.Join(staging, source, "file [5] #$.txt"), []byte("iteration five\n"))
		archiveBase := filepath.Join(root, "archive [5] #$")
		createTarXz(t, staging, archiveBase, source)
		extract := filepath.Join(root, "extract [5] #$")
		temporaryRoot := filepath.Join(root, "tmp root")
		mustMkdirAll(t, extract)
		mustMkdirAll(t, temporaryRoot)
		runWithEnvironment(t, root, append(os.Environ(), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, archiveBase)
		got, err := os.ReadFile(filepath.Join(extract, source, "file [5] #$.txt"))
		if err != nil || string(got) != "iteration five\n" {
			t.Fatalf("extracted content = %q, %v", got, err)
		}
		if info, err := os.Stat(filepath.Join(extract, source, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("archive beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createTarXz(t, root, filepath.Join(root, "-archive"), "file.txt")
		extract := filepath.Join(root, "extract")
		mustMkdirAll(t, extract)
		run(t, root, "sh", scriptFile, extract, "-archive")
		got, err := os.ReadFile(filepath.Join(extract, "file.txt"))
		if err != nil || string(got) != "content\n" {
			t.Fatalf("hyphen archive extraction = %q, %v", got, err)
		}
	})

	t.Run("xz failure prevents tar and cleans up", func(t *testing.T) {
		root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nexit 51\n")
		marker := filepath.Join(root, "tar-ran")
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\n: > \"$SHELLMAN_TAR_MARKER\"\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_TAR_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 51 {
			t.Fatalf("exit status = %d, want 51: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("tar ran after xz failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("tar failure propagates and cleans up", func(t *testing.T) {
		root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nprintf archive\n")
		mustExecutable(t, filepath.Join(fakeBin, "tar"), "#!/bin/sh\nexit 53\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 53 {
			t.Fatalf("exit status = %d, want 53: %v", exitCode(err), err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	t.Run("mktemp failure prevents xz", func(t *testing.T) {
		root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
		mustExecutable(t, filepath.Join(fakeBin, "mktemp"), "#!/bin/sh\nexit 55\n")
		marker := filepath.Join(root, "xz-ran")
		mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\n: > \"$SHELLMAN_XZ_MARKER\"\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot, "SHELLMAN_XZ_MARKER="+marker)
		_, err := runCommand(root, environment, "sh", scriptFile, extract, filepath.Join(root, "archive"))
		if exitCode(err) != 55 {
			t.Fatalf("exit status = %d, want 55: %v", exitCode(err), err)
		}
		if _, err := os.Stat(marker); !os.IsNotExist(err) {
			t.Fatalf("xz ran after mktemp failure: %v", err)
		}
		assertDirectoryEmpty(t, temporaryRoot)
	})

	for _, test := range []struct {
		name, signal string
		status       int
	}{{"INT cleanup", "-INT", 130}, {"TERM cleanup", "-TERM", 143}} {
		t.Run(test.name, func(t *testing.T) {
			root, fakeBin, extract, temporaryRoot := decompressionFailurePaths(t)
			mustExecutable(t, filepath.Join(fakeBin, "xz"), "#!/bin/sh\nkill "+test.signal+" \"$PPID\"\n")
			_, err := runCommand(root, append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"), "TMPDIR="+temporaryRoot), "sh", scriptFile, extract, filepath.Join(root, "archive"))
			if exitCode(err) != test.status {
				t.Fatalf("exit status = %d, want %d: %v", exitCode(err), test.status, err)
			}
			assertDirectoryEmpty(t, temporaryRoot)
		})
	}
}

func runnableGeneratedDecompressTarXz(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration5MigrationName].Body.(string)
	if !ok || body != currentDecompressTarXzBody {
		t.Fatal("generated iteration-5 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1|/extract/to/path, "${extractToPath}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/path/to/archive, "${pathToArchive}"|}`, `$2`)
}

func decompressionFailurePaths(t *testing.T) (root, fakeBin, extract, temporaryRoot string) {
	t.Helper()
	root = t.TempDir()
	fakeBin = filepath.Join(root, "bin")
	extract = filepath.Join(root, "extract")
	temporaryRoot = filepath.Join(root, "tmp")
	mustMkdirAll(t, fakeBin)
	mustMkdirAll(t, extract)
	mustMkdirAll(t, temporaryRoot)
	return
}

func createTarXz(t *testing.T, dir, archiveBase, source string) {
	t.Helper()
	run(t, dir, "tar", "-cf", archiveBase+".tar", source)
	run(t, dir, "xz", "-f", archiveBase+".tar")
}

func TestDecompressUnzipPlaceholderContractAndOptionSafety(t *testing.T) {
	body := currentDecompressUnzipBody
	for _, placeholder := range []string{
		`${1|/path/to/archive, "${pathToArchive}"|}`,
		`${2|/extract/to/path,"${extractToPath}"|}`,
	} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `unzip -q "${archive_path}.zip" -d "${extract_path}"`) {
		t.Fatal("snippet must quote archive and extraction paths")
	}
	if strings.Contains(body, " --") {
		t.Fatal("snippet must not rely on non-portable end-of-options syntax")
	}
}

func TestDecompressUnzipBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "zip")
	requireCommand(t, "unzip")
	scriptFile := scriptPath(t, runnableGeneratedDecompressUnzip(t))
	run(t, ".", "sh", "-n", scriptFile)

	t.Run("integrity spaces special characters and empty directory", func(t *testing.T) {
		root := t.TempDir()
		staging := filepath.Join(root, "staging")
		source := "source dir [6] #$"
		mustMkdirAll(t, filepath.Join(staging, source, "empty dir"))
		mustWriteFile(t, filepath.Join(staging, source, "file [6] #$.txt"), []byte("iteration six\n"))
		archiveBase := filepath.Join(root, "archive [6] #$")
		createZip(t, staging, archiveBase, source)
		extract := filepath.Join(root, "extract dir [6] #$")
		run(t, root, "sh", scriptFile, archiveBase, extract)
		got, err := os.ReadFile(filepath.Join(extract, source, "file [6] #$.txt"))
		if err != nil || string(got) != "iteration six\n" {
			t.Fatalf("extracted content = %q, %v", got, err)
		}
		if info, err := os.Stat(filepath.Join(extract, source, "empty dir")); err != nil || !info.IsDir() {
			t.Fatalf("empty directory was not preserved: %v", err)
		}
	})

	t.Run("archive beginning with hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		createZip(t, root, filepath.Join(root, "-archive"), "file.txt")
		extract := filepath.Join(root, "extract")
		run(t, root, "sh", scriptFile, "-archive", extract)
		got, err := os.ReadFile(filepath.Join(extract, "file.txt"))
		if err != nil || string(got) != "content\n" {
			t.Fatalf("hyphen archive extraction = %q, %v", got, err)
		}
	})

	t.Run("missing archive fails", func(t *testing.T) {
		root := t.TempDir()
		extract := filepath.Join(root, "extract")
		_, err := runCommand(root, nil, "sh", scriptFile, filepath.Join(root, "missing"), extract)
		if err == nil {
			t.Fatal("missing archive unexpectedly succeeded")
		}
		if _, err := os.Stat(extract); !os.IsNotExist(err) {
			t.Fatalf("extraction directory exists after missing archive: %v", err)
		}
	})

	t.Run("unzip failure status propagates", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "unzip"), "#!/bin/sh\nexit 61\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		_, err := runCommand(root, environment, "sh", scriptFile, filepath.Join(root, "archive"), filepath.Join(root, "extract"))
		if exitCode(err) != 61 {
			t.Fatalf("exit status = %d, want 61: %v", exitCode(err), err)
		}
	})

	t.Run("unwritable extraction target fails", func(t *testing.T) {
		root := t.TempDir()
		mustWriteFile(t, filepath.Join(root, "file.txt"), []byte("content\n"))
		archiveBase := filepath.Join(root, "archive")
		createZip(t, root, archiveBase, "file.txt")
		notDirectory := filepath.Join(root, "not-directory")
		mustWriteFile(t, notDirectory, []byte("occupied\n"))
		_, err := runCommand(root, nil, "sh", scriptFile, archiveBase, notDirectory)
		if err == nil {
			t.Fatal("file extraction target unexpectedly succeeded")
		}
	})
}

func runnableGeneratedDecompressUnzip(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration6MigrationName].Body.(string)
	if !ok || body != currentDecompressUnzipBody {
		t.Fatal("generated iteration-6 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1|/path/to/archive, "${pathToArchive}"|}`, `$1`)
	return strings.ReplaceAll(body, `${2|/extract/to/path,"${extractToPath}"|}`, `$2`)
}

func createZip(t *testing.T, dir, archiveBase, source string) {
	t.Helper()
	run(t, dir, "zip", "-rq", archiveBase+".zip", source)
}

func createTarGz(t *testing.T, dir, archiveBase, source string) {
	t.Helper()
	run(t, dir, "tar", "-cf", archiveBase+".tar", source)
	run(t, dir, "gzip", "-f", archiveBase+".tar")
}

func TestArrayFilterPlaceholderContractAndProducerStatus(t *testing.T) {
	body := currentArrayFilterBody
	for placeholder, count := range map[string]int{
		`${1:filtered}`: 1,
		`${2:myArray}`:  1,
		`${3|',"|}`:     1,
		`${3}`:          1,
		`${4:pattern}`:  1,
	} {
		if got := strings.Count(body, placeholder); got != count {
			t.Fatalf("placeholder %q occurs %d times, want %d", placeholder, got, count)
		}
	}
	if !strings.Contains(body, `grep -e ${3|',"|}${4:pattern}${3}`) {
		t.Fatal("snippet must use -e before the unchanged grep-pattern placeholders")
	}
	if !strings.HasSuffix(body, `) && wait "$!"`+"\n") {
		t.Fatal("snippet must wait for the process-substitution producer")
	}
}

func TestArrayFilterBehavior(t *testing.T) {
	requireCommand(t, "bash")
	requireCommand(t, "grep")
	body := runnableGeneratedArrayFilter(t)

	runCase := func(t *testing.T, setup, pattern string) (string, int) {
		t.Helper()
		path := filepath.Join(t.TempDir(), "array-filter.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+"status=$?\ndeclare -p filtered\nexit \"$status\"\n"))
		output, err := runCommand(".", nil, "bash", path, pattern)
		return string(output), exitCode(err)
	}

	t.Run("invalid grep regex fails", func(t *testing.T) {
		output, status := runCase(t, `source=(alpha beta); filtered=(old)`, `[`)
		if status == 0 {
			t.Fatalf("invalid regular expression unexpectedly succeeded:\n%s", output)
		}
		if !strings.Contains(output, "Invalid regular expression") {
			t.Fatalf("grep diagnostic missing:\n%s", output)
		}
	})

	t.Run("ordinary no-match succeeds", func(t *testing.T) {
		output, status := runCase(t, `source=(alpha beta); filtered=(old)`, `zzz`)
		if status != 0 {
			t.Fatalf("no-match status = %d, want 0:\n%s", status, output)
		}
		if !strings.Contains(output, `declare -a filtered=()`) {
			t.Fatalf("no-match destination was not replaced with an empty array:\n%s", output)
		}
	})

	t.Run("successful match preserves order", func(t *testing.T) {
		output, status := runCase(t, `source=(alpha beta gamma delta); filtered=(old)`, `a$`)
		if status != 0 {
			t.Fatalf("successful match status = %d, want 0:\n%s", status, output)
		}
		if !strings.Contains(output, `filtered=([0]="alpha" [1]="beta" [2]="gamma" [3]="delta")`) {
			t.Fatalf("matched destination differs or is out of order:\n%s", output)
		}
	})

	t.Run("destination assignment failure fails", func(t *testing.T) {
		output, status := runCase(t, `source=(alpha); filtered=(old); readonly -a filtered`, `alpha`)
		if status == 0 {
			t.Fatalf("readonly destination unexpectedly succeeded:\n%s", output)
		}
		if !strings.Contains(output, `filtered=([0]="old")`) {
			t.Fatalf("readonly destination changed:\n%s", output)
		}
	})
}

func runnableGeneratedArrayFilter(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration14MigrationName].Body.(string)
	if !ok || body != currentArrayFilterBody {
		t.Fatal("generated iteration-14 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:filtered}`, `filtered`)
	body = strings.ReplaceAll(body, `${2:myArray}`, `source`)
	body = strings.ReplaceAll(body, `${3|',"|}`, `"`)
	body = strings.ReplaceAll(body, `${3}`, `"`)
	body = strings.ReplaceAll(body, `${4:pattern}`, `$1`)
	return strings.ReplaceAll(body, `\${i\}`, `${i}`)
}

func TestArrayIteratePlaceholderContract(t *testing.T) {
	if got, want := len(currentArrayIterateBody), 4; got != want {
		t.Fatalf("body line count = %d, want %d", got, want)
	}
	body := strings.Join(anyStrings(t, currentArrayIterateBody), "\n")
	for _, placeholder := range []string{`${1:myArray}`, `${2:echo "\${item\}"}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays are not specified by POSIX sh.\n") {
		t.Fatal("snippet must document its indexed-array Bash dependency")
	}
}

func TestArrayIterateBehavior(t *testing.T) {
	requireCommand(t, "bash")
	body := runnableGeneratedArrayIterate(t, `printf '<%s>\n' "${item}"`)

	runCase := func(t *testing.T, setup string) (string, int) {
		t.Helper()
		path := filepath.Join(t.TempDir(), "array-iterate.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("empty array performs no action", func(t *testing.T) {
		output, status := runCase(t, `source=()`)
		if status != 0 || output != "" {
			t.Fatalf("empty iteration = status %d, output %q; want status 0 and no output", status, output)
		}
	})

	t.Run("elements retain boundaries and order", func(t *testing.T) {
		output, status := runCase(t, `source=('two words' '*' '' '-n')`)
		if status != 0 {
			t.Fatalf("iteration status = %d, want 0: %s", status, output)
		}
		if want := "<two words>\n<*>\n<>\n<-n>\n"; output != want {
			t.Fatalf("iteration output = %q, want %q", output, want)
		}
	})

	t.Run("single action failure propagates", func(t *testing.T) {
		failureBody := runnableGeneratedArrayIterate(t, `false`)
		path := filepath.Join(t.TempDir(), "array-iterate-failure.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\nsource=(one)\n"+failureBody))
		_, err := runCommand(".", nil, "bash", path)
		if exitCode(err) == 0 {
			t.Fatal("failing action unexpectedly succeeded")
		}
	})
}

func runnableGeneratedArrayIterate(t *testing.T, action string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration15MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentArrayIterateBody) {
		t.Fatal("generated iteration-15 body differs from candidate")
	}
	lines := anyStrings(t, body)
	result := strings.Join(lines, "\n")
	result = strings.ReplaceAll(result, `${1:myArray}`, `source`)
	return strings.ReplaceAll(result, `${2:echo "\${item\}"}`, action)
}

func anyStrings(t *testing.T, values []any) []string {
	t.Helper()
	stringsOnly := make([]string, len(values))
	for index, value := range values {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d is %T, want string", index, value)
		}
		stringsOnly[index] = line
	}
	return stringsOnly
}

func TestArrayLengthPlaceholderContract(t *testing.T) {
	body := currentArrayLengthBody
	for _, placeholder := range []string{`${1:length}`, `${2:myArray}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `${1:length}=${#${2:myArray}[@]}`) {
		t.Fatal("destination and nested array-length placeholders changed")
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays are not specified by POSIX sh.\n") {
		t.Fatal("snippet must document its indexed-array Bash dependency")
	}
}

func TestArrayLengthBehavior(t *testing.T) {
	requireCommand(t, "bash")
	body := runnableGeneratedArrayLength(t)

	runCase := func(t *testing.T, setup string) (string, int) {
		t.Helper()
		path := filepath.Join(t.TempDir(), "array-length.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+"status=$?\nprintf '%s:%s\\n' \"$status\" \"${length-UNSET}\"\nexit \"$status\"\n"))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("empty array has length zero", func(t *testing.T) {
		output, status := runCase(t, `source=(); length=old`)
		if status != 0 || output != "0:0\n" {
			t.Fatalf("empty length = status %d, output %q; want status 0, output %q", status, output, "0:0\n")
		}
	})

	t.Run("counts assigned sparse elements", func(t *testing.T) {
		output, status := runCase(t, `source=(); source[2]='two words'; source[9]='*'; source[20]=''`)
		if status != 0 || output != "0:3\n" {
			t.Fatalf("sparse length = status %d, output %q; want status 0, output %q", status, output, "0:3\n")
		}
	})

	t.Run("unset array has length zero", func(t *testing.T) {
		output, status := runCase(t, `unset source`)
		if status != 0 || output != "0:0\n" {
			t.Fatalf("unset length = status %d, output %q; want status 0, output %q", status, output, "0:0\n")
		}
	})

	t.Run("readonly destination failure propagates", func(t *testing.T) {
		output, status := runCase(t, `source=(one two); readonly length=old`)
		if status == 0 {
			t.Fatalf("readonly destination unexpectedly succeeded: %s", output)
		}
		if !strings.Contains(output, "readonly variable") {
			t.Fatalf("readonly assignment diagnostic missing: %s", output)
		}
	})
}

func runnableGeneratedArrayLength(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration16MigrationName].Body.(string)
	if !ok || body != currentArrayLengthBody {
		t.Fatal("generated iteration-16 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:length}`, `length`)
	return strings.ReplaceAll(body, `${2:myArray}`, `source`)
}

func TestArrayPrintPlaceholderContractAndSafety(t *testing.T) {
	body := currentArrayPrintBody
	if strings.Count(body, `${1:myArray}`) != 1 {
		t.Fatal("array placeholder must occur exactly once")
	}
	if !strings.Contains(body, `(IFS=' '; printf '%s\n' "\${${1:myArray}[*]}")`) {
		t.Fatal("snippet must use a fixed-space join and printf")
	}
	if strings.Contains(body, "echo ") {
		t.Fatal("snippet must not use implementation-dependent echo for data")
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays are not specified by POSIX sh.\n") {
		t.Fatal("snippet must document its indexed-array Bash dependency")
	}
}

func TestArrayPrintBehavior(t *testing.T) {
	requireCommand(t, "bash")
	body := runnableGeneratedArrayPrint(t)

	runCase := func(t *testing.T, setup, trailer string) (string, int) {
		t.Helper()
		path := filepath.Join(t.TempDir(), "array-print.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("empty array prints newline", func(t *testing.T) {
		output, status := runCase(t, `source=()`, "")
		if status != 0 || output != "\n" {
			t.Fatalf("empty print = status %d, output %q; want status 0, newline", status, output)
		}
	})

	t.Run("joins elements safely with spaces", func(t *testing.T) {
		output, status := runCase(t, `source=('two words' '*' '' '-n' '\c')`, "")
		if status != 0 {
			t.Fatalf("array print status = %d, want 0: %s", status, output)
		}
		if want := "two words *  -n \\c\n"; output != want {
			t.Fatalf("array print output = %q, want %q", output, want)
		}
	})

	t.Run("caller IFS is ignored and preserved", func(t *testing.T) {
		output, status := runCase(t, "IFS=,; source=(one two three)", `printf '<%s>\n' "$IFS"`+"\n")
		if status != 0 || output != "one two three\n<,>\n" {
			t.Fatalf("custom-IFS print = status %d, output %q", status, output)
		}
	})

	t.Run("printf write failure propagates", func(t *testing.T) {
		path := filepath.Join(t.TempDir(), "array-print-failure.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\nsource=(data)\nexec >/dev/full\n"+body+"status=$?\nexit \"$status\"\n"))
		output, err := runCommand(".", nil, "bash", path)
		if exitCode(err) == 0 {
			t.Fatalf("write failure unexpectedly succeeded: %s", output)
		}
	})
}

func runnableGeneratedArrayPrint(t *testing.T) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration17MigrationName].Body.(string)
	if !ok || body != currentArrayPrintBody {
		t.Fatal("generated iteration-17 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:myArray}`, `source`)
	return strings.ReplaceAll(body, `\${source[*]}`, `${source[*]}`)
}

func TestArrayPushPlaceholderContract(t *testing.T) {
	body := currentArrayPushBody
	for _, placeholder := range []string{`${1:myArray}`, `${2:newItem}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `${1:myArray}+=('${2:newItem}')`) {
		t.Fatal("append assignment and placeholder order changed")
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays and compound assignment are not specified by POSIX sh.\n") {
		t.Fatal("snippet must document its Bash array and compound-assignment dependencies")
	}
}

func TestArrayPushBehavior(t *testing.T) {
	requireCommand(t, "bash")

	runCase := func(t *testing.T, setup, item, assertion string) (string, int) {
		t.Helper()
		body := runnableGeneratedArrayPush(t, item)
		path := filepath.Join(t.TempDir(), "array-push.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+assertion+"\n"))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("appends one element preserving boundaries", func(t *testing.T) {
		item := "two words * -n\nsecond line $HOME"
		assertion := `[[ ${#source[@]} -eq 2 && ${source[0]} == old && ${source[1]} == $'two words * -n\nsecond line $HOME' ]]`
		output, status := runCase(t, `source=(old)`, item, assertion)
		if status != 0 {
			t.Fatalf("boundary-preserving append failed with status %d: %s", status, output)
		}
	})

	t.Run("empty element is retained", func(t *testing.T) {
		assertion := `[[ ${#source[@]} -eq 1 && -v 'source[0]' && -z ${source[0]} ]]`
		output, status := runCase(t, `source=()`, "", assertion)
		if status != 0 {
			t.Fatalf("empty append failed with status %d: %s", status, output)
		}
	})

	t.Run("sparse array appends after highest index", func(t *testing.T) {
		assertion := `[[ ${#source[@]} -eq 2 && ${source[4]} == old && ${source[5]} == new ]]`
		output, status := runCase(t, `source=(); source[4]=old`, "new", assertion)
		if status != 0 {
			t.Fatalf("sparse append failed with status %d: %s", status, output)
		}
	})

	t.Run("readonly array failure propagates", func(t *testing.T) {
		output, status := runCase(t, `source=(old); readonly -a source`, "new", `status=$?; exit "$status"`)
		if status == 0 {
			t.Fatalf("readonly append unexpectedly succeeded: %s", output)
		}
		if !strings.Contains(output, "readonly variable") {
			t.Fatalf("readonly diagnostic missing: %s", output)
		}
	})
}

func runnableGeneratedArrayPush(t *testing.T, item string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration18MigrationName].Body.(string)
	if !ok || body != currentArrayPushBody {
		t.Fatal("generated iteration-18 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:myArray}`, `source`)
	return strings.ReplaceAll(body, `${2:newItem}`, item)
}

func TestArrayRangePlaceholderContractAndBoundaries(t *testing.T) {
	body := currentArrayRangeBody
	for _, placeholder := range []string{`${1:newArray}`, `${2:myArray}`, `${3:fromIndex}`, `${4:n}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `${1:newArray}=("${${2:myArray}[@]:${3:fromIndex}:${4:n}}")`) {
		t.Fatal("array-assignment slice or placeholder nesting changed")
	}
	if strings.Contains(body, `[*]`) {
		t.Fatal("slice must not collapse element boundaries")
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays and array slicing are not specified by POSIX sh.\n") {
		t.Fatal("snippet must document its Bash array and slicing dependencies")
	}
}

func TestArrayRangeBehavior(t *testing.T) {
	requireCommand(t, "bash")

	runCase := func(t *testing.T, setup, destination, source, from, count, assertion string) (string, int) {
		t.Helper()
		body := runnableGeneratedArrayRange(t, destination, source, from, count)
		path := filepath.Join(t.TempDir(), "array-range.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+assertion+"\n"))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("slice preserves boundaries and order", func(t *testing.T) {
		setup := `source=(zero 'two words' '' '*' '-n' $'line one\nline two' six)`
		assertion := `[[ ${#result[@]} -eq 5 && ${result[0]} == 'two words' && -z ${result[1]} && ${result[2]} == '*' && ${result[3]} == '-n' && ${result[4]} == $'line one\nline two' ]]`
		output, status := runCase(t, setup, "result", "source", "1", "5", assertion)
		if status != 0 {
			t.Fatalf("boundary-preserving slice failed with status %d: %s", status, output)
		}
	})

	t.Run("empty source replaces destination with empty array", func(t *testing.T) {
		assertion := `[[ ${#result[@]} -eq 0 ]]`
		output, status := runCase(t, `source=(); result=(old)`, "result", "source", "0", "3", assertion)
		if status != 0 {
			t.Fatalf("empty slice failed with status %d: %s", status, output)
		}
	})

	t.Run("sparse array slices assigned-element order", func(t *testing.T) {
		setup := `source=(); source[2]=two; source[8]=eight; source[20]=twenty`
		assertion := `[[ ${#result[@]} -eq 2 && ${result[0]} == eight && ${result[1]} == twenty ]]`
		output, status := runCase(t, setup, "result", "source", "3", "2", assertion)
		if status != 0 {
			t.Fatalf("sparse slice failed with status %d: %s", status, output)
		}
	})

	t.Run("source may also be destination", func(t *testing.T) {
		assertion := `[[ ${#source[@]} -eq 2 && ${source[0]} == one && ${source[1]} == two ]]`
		output, status := runCase(t, `source=(zero one two three)`, "source", "source", "1", "2", assertion)
		if status != 0 {
			t.Fatalf("in-place slice failed with status %d: %s", status, output)
		}
	})

	t.Run("readonly destination failure propagates", func(t *testing.T) {
		assertion := `status=$?; exit "$status"`
		output, status := runCase(t, `source=(zero one); result=(old); readonly -a result`, "result", "source", "0", "1", assertion)
		if status == 0 {
			t.Fatalf("readonly slice unexpectedly succeeded: %s", output)
		}
		if !strings.Contains(output, "readonly variable") {
			t.Fatalf("readonly diagnostic missing: %s", output)
		}
	})

	t.Run("negative length fails", func(t *testing.T) {
		output, status := runCase(t, `source=(zero one); result=(old)`, "result", "source", "0", "-1", `status=$?; exit "$status"`)
		if status == 0 {
			t.Fatalf("negative slice length unexpectedly succeeded: %s", output)
		}
	})
}

func runnableGeneratedArrayRange(t *testing.T, destination, source, from, count string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration19MigrationName].Body.(string)
	if !ok || body != currentArrayRangeBody {
		t.Fatal("generated iteration-19 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:newArray}`, destination)
	body = strings.ReplaceAll(body, `${2:myArray}`, source)
	body = strings.ReplaceAll(body, `${3:fromIndex}`, from)
	return strings.ReplaceAll(body, `${4:n}`, count)
}

func TestArrayReplacePlaceholderContractAndBoundaries(t *testing.T) {
	body := currentArrayReplaceBody
	for _, placeholder := range []string{`${1:newArray}`, `${2:myArray}`, `${3:find}`, `${4:replace}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q must occur exactly once", placeholder)
		}
	}
	if !strings.Contains(body, `${1:newArray}=("${${2:myArray}[@]//${3:find}/${4:replace}}")`) {
		t.Fatal("array substitution assignment or placeholder nesting changed")
	}
	if strings.Contains(body, `[*]`) {
		t.Fatal("replacement must not collapse element boundaries")
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays and pattern substitution are not specified by POSIX sh.\n") {
		t.Fatal("snippet must document its Bash array and pattern-substitution dependencies")
	}
}

func TestArrayReplaceBehavior(t *testing.T) {
	requireCommand(t, "bash")

	runCase := func(t *testing.T, setup, destination, source, find, replacement, assertion string) (string, int) {
		t.Helper()
		body := runnableGeneratedArrayReplace(t, destination, source, find, replacement)
		path := filepath.Join(t.TempDir(), "array-replace.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+assertion+"\n"))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("replacement preserves boundaries and order", func(t *testing.T) {
		setup := `source=('foo bar' '' '*' '-foo' $'foo\nbar' '$HOME')`
		assertion := `[[ ${#result[@]} -eq 6 && ${result[0]} == 'X bar' && -z ${result[1]} && ${result[2]} == '*' && ${result[3]} == '-X' && ${result[4]} == $'X\nbar' && ${result[5]} == '$HOME' ]]`
		output, status := runCase(t, setup, "result", "source", "foo", "X", assertion)
		if status != 0 {
			t.Fatalf("boundary-preserving replacement failed with status %d: %s", status, output)
		}
	})

	t.Run("find uses Bash patterns", func(t *testing.T) {
		assertion := `[[ ${#result[@]} -eq 2 && ${result[0]} == 'v#.#' && ${result[1]} == 'no digits' ]]`
		output, status := runCase(t, `source=('v1.2' 'no digits')`, "result", "source", `[[:digit:]]`, `#`, assertion)
		if status != 0 {
			t.Fatalf("pattern replacement failed with status %d: %s", status, output)
		}
	})

	t.Run("empty source replaces destination with empty array", func(t *testing.T) {
		output, status := runCase(t, `source=(); result=(old)`, "result", "source", "x", "y", `[[ ${#result[@]} -eq 0 ]]`)
		if status != 0 {
			t.Fatalf("empty replacement failed with status %d: %s", status, output)
		}
	})

	t.Run("sparse source retains assigned order", func(t *testing.T) {
		setup := `source=(); source[3]=cat; source[9]='cat dog'`
		assertion := `[[ ${#result[@]} -eq 2 && ${result[0]} == fox && ${result[1]} == 'fox dog' ]]`
		output, status := runCase(t, setup, "result", "source", "cat", "fox", assertion)
		if status != 0 {
			t.Fatalf("sparse replacement failed with status %d: %s", status, output)
		}
	})

	t.Run("source may also be destination", func(t *testing.T) {
		assertion := `[[ ${#source[@]} -eq 2 && ${source[0]} == X && ${source[1]} == bX ]]`
		output, status := runCase(t, `source=(a ba)`, "source", "source", "a", "X", assertion)
		if status != 0 {
			t.Fatalf("in-place replacement failed with status %d: %s", status, output)
		}
	})

	t.Run("readonly destination failure propagates", func(t *testing.T) {
		assertion := `status=$?; exit "$status"`
		output, status := runCase(t, `source=(cat); result=(old); readonly -a result`, "result", "source", "cat", "fox", assertion)
		if status == 0 {
			t.Fatalf("readonly replacement unexpectedly succeeded: %s", output)
		}
		if !strings.Contains(output, "readonly variable") {
			t.Fatalf("readonly diagnostic missing: %s", output)
		}
	})
}

func runnableGeneratedArrayReplace(t *testing.T, destination, source, find, replacement string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration20MigrationName].Body.(string)
	if !ok || body != currentArrayReplaceBody {
		t.Fatal("generated iteration-20 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:newArray}`, destination)
	body = strings.ReplaceAll(body, `${2:myArray}`, source)
	body = strings.ReplaceAll(body, `${3:find}`, find)
	return strings.ReplaceAll(body, `${4:replace}`, replacement)
}

func TestArrayReversePlaceholderContractAndBashDocumentation(t *testing.T) {
	body := strings.Join(anyStrings(t, currentArrayReverseBody), "\n")
	if !strings.Contains(body, `${1:myArray}`) || !strings.Contains(body, `${2:reversed}`) {
		t.Fatal("the original source and destination placeholders must both remain")
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays and C-style for loops are not specified by POSIX sh.\n") {
		t.Fatal("Bash-only comment must identify the actual shell dependencies")
	}
	if strings.Contains(body, "echo ") || strings.Contains(body, "+=") {
		t.Fatal("reverse must neither use echo nor append to the destination")
	}
}

func TestArrayReverseBehavior(t *testing.T) {
	requireCommand(t, "bash")

	runCase := func(t *testing.T, setup, destination, source, trailer string) (string, int) {
		t.Helper()
		body := runnableGeneratedArrayReverse(t, destination, source)
		path := filepath.Join(t.TempDir(), "array-reverse.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	t.Run("dense array reverses element order and preserves values", func(t *testing.T) {
		setup := `source=('first value' '' '*' '-n' 'a\\b' $'line one\nline two' '$HOME' '"quoted"')`
		trailer := `[[ $snippet_status -eq 0 && ${#reversed[@]} -eq 8 && ${reversed[0]} == '"quoted"' && ${reversed[1]} == '$HOME' && ${reversed[2]} == $'line one\nline two' && ${reversed[3]} == 'a\\b' && ${reversed[4]} == '-n' && ${reversed[5]} == '*' && -z ${reversed[6]} && ${reversed[7]} == 'first value' && ! -v source ]]`
		output, status := runCase(t, setup, "reversed", "source", trailer)
		if status != 0 {
			t.Fatalf("dense reverse failed with status %d: %s", status, output)
		}
	})

	t.Run("empty array replaces destination", func(t *testing.T) {
		output, status := runCase(t, `source=(); reversed=(old)`, "reversed", "source", `[[ $snippet_status -eq 0 && ${#reversed[@]} -eq 0 && ! -v source ]]`)
		if status != 0 {
			t.Fatalf("empty reverse failed with status %d: %s", status, output)
		}
	})

	t.Run("single element remains one element", func(t *testing.T) {
		output, status := runCase(t, `source=('one item'); reversed=(old)`, "reversed", "source", `[[ $snippet_status -eq 0 && ${#reversed[@]} -eq 1 && ${reversed[0]} == 'one item' ]]`)
		if status != 0 {
			t.Fatalf("single-element reverse failed with status %d: %s", status, output)
		}
	})

	t.Run("preexisting destination is replaced", func(t *testing.T) {
		output, status := runCase(t, `source=(a b); reversed=(old stale)`, "reversed", "source", `[[ $snippet_status -eq 0 && ${#reversed[@]} -eq 2 && ${reversed[0]} == b && ${reversed[1]} == a ]]`)
		if status != 0 {
			t.Fatalf("replacement failed: status %d: %s", status, output)
		}
	})

	t.Run("sparse source reverses assigned element order", func(t *testing.T) {
		output, status := runCase(t, `source=(); source[2]='two value'; source[8]=''; source[21]='twenty one'; reversed=(old)`, "reversed", "source", `[[ $snippet_status -eq 0 && ${#reversed[@]} -eq 3 && ${reversed[0]} == 'twenty one' && -z ${reversed[1]} && ${reversed[2]} == 'two value' ]]`)
		if status != 0 {
			t.Fatalf("sparse reverse failed: status %d: %s", status, output)
		}
	})

	t.Run("readonly empty destination replacement failure propagates", func(t *testing.T) {
		output, status := runCase(t, `source=(one two); reversed=(); readonly -a reversed`, "reversed", "source", `exit "$snippet_status"`)
		if status == 0 {
			t.Fatalf("readonly assignment was masked: %q", output)
		}
		if !strings.Contains(output, "readonly variable") {
			t.Fatalf("readonly failure diagnostic not observed: %q", output)
		}
	})

	t.Run("readonly populated destination reset failure propagates", func(t *testing.T) {
		output, status := runCase(t, `source=(one two); reversed=(old); readonly -a reversed`, "reversed", "source", `exit "$snippet_status"`)
		if status == 0 || !strings.Contains(output, "readonly variable") {
			t.Fatalf("readonly reset failure was masked: status=%d output=%q", status, output)
		}
	})

	t.Run("source and destination may be the same", func(t *testing.T) {
		output, status := runCase(t, `source=(); source[3]=three; source[9]=nine`, "source", "source", `[[ $snippet_status -eq 0 && ${#source[@]} -eq 2 && ${source[0]} == nine && ${source[1]} == three ]]`)
		if status != 0 {
			t.Fatalf("same-name reverse failed: status=%d output=%q", status, output)
		}
	})

	t.Run("optional unset failure propagates", func(t *testing.T) {
		output, status := runCase(t, `source=(one two); readonly -a source`, "reversed", "source", `exit "$snippet_status"`)
		if status == 0 || !strings.Contains(output, "readonly variable") {
			t.Fatalf("optional unset failure was masked: status=%d output=%q", status, output)
		}
	})

	t.Run("printf does not interpret leading hyphens or backslashes", func(t *testing.T) {
		output, status := runCase(t, `source=('a\b' '-n')`, "reversed", "source", `[[ $snippet_status -eq 0 ]]`)
		if status != 0 || output != "-n a\\b\n" {
			t.Fatalf("safe output failed: status=%d output=%q", status, output)
		}
	})

	t.Run("output write failure propagates", func(t *testing.T) {
		body := runnableGeneratedArrayReverse(t, "reversed", "source")
		path := filepath.Join(t.TempDir(), "array-reverse-output-failure.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\nsource=(one two)\n"+body))
		_, err := runCommand(".", nil, "bash", "-c", `bash "$1" >/dev/full`, "bash", path)
		if exitCode(err) == 0 {
			t.Fatal("printf write failure was masked")
		}
	})
}

func runnableGeneratedArrayReverse(t *testing.T, destination, source string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration21MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentArrayReverseBody) {
		t.Fatal("generated iteration-21 body differs from candidate")
	}
	result := strings.Join(anyStrings(t, body), "\n")
	result = strings.ReplaceAll(result, `${1:myArray}`, source)
	result = strings.ReplaceAll(result, `${2:reversed}`, destination)
	return strings.ReplaceAll(result, `\${`, `${`)
}

func TestCommandFailureCheckPlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := strings.Join(anyStrings(t, currentCommandFailureCheckBody), "\n")
	if strings.Count(body, `${1:command}`) != 1 {
		t.Fatalf("command placeholder count = %d, want 1", strings.Count(body, `${1:command}`))
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("failure-check must remain POSIX sh")
	}
	if !strings.Contains(body, `(exit "\${_shellman_command_status}")`) {
		t.Fatal("failure-check must return the checked command's captured status")
	}
}

func TestCommandFailureCheckBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, command, trailer string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandFailureCheck(t, command)
		path := filepath.Join(t.TempDir(), "command-failure-check.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("success reports success and returns zero", func(t *testing.T) {
		output, status := runCase(t, "", ":", `exit "$snippet_status"`)
		if status != 0 || output != "succeed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("status one reports failure and returns one", func(t *testing.T) {
		output, status := runCase(t, "", "false", `exit "$snippet_status"`)
		if status != 1 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("arbitrary nonzero status is preserved", func(t *testing.T) {
		output, status := runCase(t, `checked() { return 42; }`, "checked", `exit "$snippet_status"`)
		if status != 42 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("command not found is quiet and preserves 127", func(t *testing.T) {
		output, status := runCase(t, "", "shellman_command_that_does_not_exist", `exit "$snippet_status"`)
		if status != 127 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("checked stdout and stderr are suppressed", func(t *testing.T) {
		setup := `checked() { printf 'command stdout\n'; printf 'command stderr\n' >&2; return 0; }`
		output, status := runCase(t, setup, "checked", `exit "$snippet_status"`)
		if status != 0 || output != "succeed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("command arguments retain caller shell semantics", func(t *testing.T) {
		setup := `checked() { [ "$1" = 'space value' ] && [ "$2" = '*' ] && [ "$3" = '-n' ] && [ "$4" = '$HOME' ]; }`
		command := `checked 'space value' '*' '-n' '$HOME'`
		output, status := runCase(t, setup, command, `exit "$snippet_status"`)
		if status != 0 || output != "succeed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("shell-function side effects remain in caller shell", func(t *testing.T) {
		setup := `marker=before; checked() { marker='after value'; return 9; }`
		output, status := runCase(t, setup, "checked", `[ "$snippet_status" -eq 9 ] && [ "$marker" = 'after value' ]`)
		if status != 0 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("status capture failure is not masked", func(t *testing.T) {
		setup := `_shellman_command_status=7; readonly _shellman_command_status`
		output, status := runCase(t, setup, ":", `exit "$snippet_status"`)
		if status == 0 || !strings.Contains(output, "readonly") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("diagnostic write failure is not masked", func(t *testing.T) {
		body := runnableGeneratedCommandFailureCheck(t, "false")
		path := filepath.Join(t.TempDir(), "command-failure-check-output.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+body))
		_, err := runCommand(".", nil, "sh", "-c", `sh "$1" >/dev/full`, "sh", path)
		if exitCode(err) == 0 {
			t.Fatal("printf write failure was masked")
		}
	})
}

func runnableGeneratedCommandFailureCheck(t *testing.T, command string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration23MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentCommandFailureCheckBody) {
		t.Fatal("generated iteration-23 body differs from candidate")
	}
	result := strings.Join(anyStrings(t, body), "\n")
	result = strings.ReplaceAll(result, `${1:command}`, command)
	return strings.ReplaceAll(result, `\${`, `${`)
}

func TestCommandHideErrorPlaceholderContractAndPOSIXSyntax(t *testing.T) {
	if strings.Count(currentCommandHideErrorBody, `${1:command}`) != 1 {
		t.Fatalf("command placeholder count = %d, want 1", strings.Count(currentCommandHideErrorBody, `${1:command}`))
	}
	if strings.Contains(currentCommandHideErrorBody, "[[") || strings.Contains(currentCommandHideErrorBody, "function ") || strings.Contains(currentCommandHideErrorBody, "local ") {
		t.Fatal("hide-error must remain POSIX sh")
	}
	if !strings.HasPrefix(currentCommandHideErrorBody, "{\n") || !strings.HasSuffix(currentCommandHideErrorBody, "\n} 2>/dev/null\n") {
		t.Fatal("stderr redirection must apply to the complete raw command placeholder")
	}
}

func TestCommandHideErrorBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, command, trailer string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandHideError(t, command)
		path := filepath.Join(t.TempDir(), "command-hide-error.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("success preserves stdout and zero status", func(t *testing.T) {
		output, status := runCase(t, "", `printf 'visible stdout\n'; printf 'hidden stderr\n' >&2`, `exit "$snippet_status"`)
		if status != 0 || output != "visible stdout\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("status one is preserved", func(t *testing.T) {
		output, status := runCase(t, "", "false", `exit "$snippet_status"`)
		if status != 1 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("arbitrary nonzero status is preserved", func(t *testing.T) {
		output, status := runCase(t, `checked() { return 42; }`, "checked", `exit "$snippet_status"`)
		if status != 42 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("command not found is quiet and preserves 127", func(t *testing.T) {
		output, status := runCase(t, "", "shellman_command_that_does_not_exist", `exit "$snippet_status"`)
		if status != 127 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("redirection covers the complete raw command list", func(t *testing.T) {
		command := `printf 'first hidden\n' >&2; printf 'second hidden\n' >&2; printf 'list stdout\n'`
		output, status := runCase(t, "", command, `exit "$snippet_status"`)
		if status != 0 || output != "list stdout\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("arguments retain caller shell semantics", func(t *testing.T) {
		setup := `checked() { [ "$1" = 'space value' ] && [ "$2" = '*' ] && [ "$3" = '-n' ] && [ "$4" = '$HOME' ] && [ "$5" = 'quote"backslash\line' ]; }`
		command := `checked 'space value' '*' '-n' '$HOME' 'quote"backslash\line'`
		output, status := runCase(t, setup, command, `exit "$snippet_status"`)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("side effects remain in caller shell", func(t *testing.T) {
		output, status := runCase(t, `marker=before`, `marker='after value'; false`, `[ "$snippet_status" -eq 1 ] && [ "$marker" = 'after value' ]`)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("pipeline status retains POSIX last-command semantics", func(t *testing.T) {
		output, status := runCase(t, "", `printf 'hidden pipeline\n' >&2 | false`, `exit "$snippet_status"`)
		if status != 1 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCommandHideError(t *testing.T, command string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration24MigrationName].Body.(string)
	if !ok || body != currentCommandHideErrorBody {
		t.Fatal("generated iteration-24 body differs from candidate")
	}
	return strings.ReplaceAll(body, `${1:command}`, command)
}

func TestCommandIfExistsPlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := strings.Join(anyStrings(t, currentCommandIfExistsBody), "\n")
	if strings.Count(body, `${1:command}`) != 2 {
		t.Fatalf("command-name placeholder count = %d, want 2", strings.Count(body, `${1:command}`))
	}
	if strings.Count(body, `${2:echo "command \"${1:command}\" exists on system"}`) != 1 {
		t.Fatal("raw action placeholder, nested command-name mirror, or escaping changed")
	}
	if !strings.Contains(body, `command -v -- "${1:command}" >/dev/null 2>&1`) {
		t.Fatal("command name must be one quoted operand and probe output must be suppressed")
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("if-exists must remain POSIX sh")
	}
}

func TestCommandIfExistsBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, commandName, action, trailer string, environment []string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandIfExists(t, commandName, action)
		path := filepath.Join(t.TempDir(), "command-if-exists.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("existing command runs raw action", func(t *testing.T) {
		output, status := runCase(t, "", "sh", `printf 'action stdout\n'; printf 'action stderr\n' >&2`, `exit "$snippet_status"`, nil)
		if status != 0 || output != "action stdout\naction stderr\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("missing command is quiet and skips action", func(t *testing.T) {
		output, status := runCase(t, "", "shellman_command_that_does_not_exist", `printf 'must not run\n'; false`, `exit "$snippet_status"`, nil)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("empty command name is quiet and skips action", func(t *testing.T) {
		output, status := runCase(t, "", "", `printf 'must not run\n'; false`, `exit "$snippet_status"`, nil)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("action arbitrary failure status propagates", func(t *testing.T) {
		setup := `fail_action() { return 42; }`
		output, status := runCase(t, setup, "sh", "fail_action", `exit "$snippet_status"`, nil)
		if status != 42 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("action status one propagates", func(t *testing.T) {
		output, status := runCase(t, "", "sh", "false", `exit "$snippet_status"`, nil)
		if status != 1 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("action command-not-found status and diagnostic propagate", func(t *testing.T) {
		output, status := runCase(t, "", "sh", "shellman_action_that_does_not_exist", `exit "$snippet_status"`, nil)
		if status != 127 || !strings.Contains(output, "shellman_action_that_does_not_exist") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("raw action pipeline uses POSIX last-command status", func(t *testing.T) {
		setup := `fail_action() { return 42; }`
		output, status := runCase(t, setup, "sh", `printf 'pipeline input\n' | fail_action`, `exit "$snippet_status"`, nil)
		if status != 42 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("raw action list and caller side effects are preserved", func(t *testing.T) {
		setup := `marker=before`
		action := `marker='after value'; printf 'one\n'; printf 'two\n'`
		trailer := `[ "$snippet_status" -eq 0 ] && [ "$marker" = 'after value' ]`
		output, status := runCase(t, setup, "sh", action, trailer, nil)
		if status != 0 || output != "one\ntwo\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("quoted runtime command-name data is one operand", func(t *testing.T) {
		fakeBin := t.TempDir()
		names := []string{"space name", "*", "-n", `$HOME`, "quote\"backslash\\name", "line\nbreak"}
		for _, name := range names {
			mustExecutable(t, filepath.Join(fakeBin, name), "#!/bin/sh\nexit 0\n")
			setup := "candidate='" + strings.ReplaceAll(name, "'", `'\\''`) + "'"
			environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
			output, status := runCase(t, setup, `${candidate}`, `:`, `exit "$snippet_status"`, environment)
			if status != 0 || output != "" {
				t.Fatalf("name=%q status=%d output=%q", name, status, output)
			}
		}
	})
}

func runnableGeneratedCommandIfExists(t *testing.T, commandName, action string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration25MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentCommandIfExistsBody) {
		t.Fatal("generated iteration-25 body differs from candidate")
	}
	result := strings.Join(anyStrings(t, body), "\n")
	result = strings.ReplaceAll(result, `${2:echo "command \"${1:command}\" exists on system"}`, action)
	return strings.ReplaceAll(result, `${1:command}`, commandName)
}

func TestCommandNicePlaceholderContractAndPOSIXSyntax(t *testing.T) {
	if strings.Count(currentCommandNiceBody, `${1|-20,-15,-10,-5,0,5,10,15,19|}`) != 1 {
		t.Fatal("niceness choice placeholder changed")
	}
	if strings.Count(currentCommandNiceBody, `${2:command}`) != 1 {
		t.Fatal("raw command-invocation placeholder changed")
	}
	if !strings.Contains(currentCommandNiceBody, `nice -n ${1|-20,-15,-10,-5,0,5,10,15,19|} -- ${2:command}`) {
		t.Fatal("nice option termination must precede the raw command invocation")
	}
	if strings.Contains(currentCommandNiceBody, "[[") || strings.Contains(currentCommandNiceBody, "function ") || strings.Contains(currentCommandNiceBody, "local ") {
		t.Fatal("nice snippet shell syntax must remain POSIX sh")
	}
}

func TestCommandNiceBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "nice")

	runCase := func(t *testing.T, fakeBin, command string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandNice(t, "5", command)
		path := filepath.Join(t.TempDir(), "command-nice.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+body))
		run(t, ".", "sh", "-n", path)
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	newFakeBin := func(t *testing.T) string {
		t.Helper()
		fakeBin := t.TempDir()
		mustExecutable(t, filepath.Join(fakeBin, "sudo"), "#!/bin/sh\nexec \"$@\"\n")
		return fakeBin
	}

	t.Run("success and streams pass through", func(t *testing.T) {
		fakeBin := newFakeBin(t)
		mustExecutable(t, filepath.Join(fakeBin, "checked"), "#!/bin/sh\nprintf 'command stdout\\n'\nprintf 'command stderr\\n' >&2\nexit 0\n")
		output, status := runCase(t, fakeBin, "checked")
		if status != 0 || output != "command stdout\ncommand stderr\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, tc := range []struct {
		name   string
		status int
	}{
		{"status one", 1},
		{"arbitrary status", 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			fakeBin := newFakeBin(t)
			mustExecutable(t, filepath.Join(fakeBin, "checked"), fmt.Sprintf("#!/bin/sh\nexit %d\n", tc.status))
			output, status := runCase(t, fakeBin, "checked")
			if status != tc.status || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("command not found status propagates", func(t *testing.T) {
		fakeBin := newFakeBin(t)
		output, status := runCase(t, fakeBin, "shellman_command_that_does_not_exist")
		if status != 127 || !strings.Contains(output, "shellman_command_that_does_not_exist") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("sudo failure status propagates", func(t *testing.T) {
		fakeBin := t.TempDir()
		mustExecutable(t, filepath.Join(fakeBin, "sudo"), "#!/bin/sh\nprintf 'sudo failed\\n' >&2\nexit 77\n")
		output, status := runCase(t, fakeBin, "true")
		if status != 77 || output != "sudo failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("empty raw invocation exposes nice no-command behavior", func(t *testing.T) {
		fakeBin := newFakeBin(t)
		output, status := runCase(t, fakeBin, "")
		if status != 125 || !strings.Contains(output, "command must be given") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("command name beginning with hyphen is not a nice option", func(t *testing.T) {
		fakeBin := newFakeBin(t)
		mustExecutable(t, filepath.Join(fakeBin, "-shellman-option"), "#!/bin/sh\nprintf 'leading hyphen ran\\n'\n")
		output, status := runCase(t, fakeBin, "-shellman-option")
		if status != 0 || output != "leading hyphen ran\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("quoted command path and arguments retain data", func(t *testing.T) {
		fakeBin := newFakeBin(t)
		commandPath := filepath.Join(fakeBin, "space command")
		mustExecutable(t, commandPath, "#!/bin/sh\nprintf '<%s>\\n' \"$@\"\n")
		command := fmt.Sprintf("'%s' 'space value' '*' '-n' '$HOME' 'quote\"backslash\\line'", commandPath)
		output, status := runCase(t, fakeBin, command)
		want := "<space value>\n<*>\n<-n>\n<$HOME>\n<quote\"backslash\\line>\n"
		if status != 0 || output != want {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("raw pipeline retains caller shell pipeline semantics", func(t *testing.T) {
		fakeBin := newFakeBin(t)
		mustExecutable(t, filepath.Join(fakeBin, "produce"), "#!/bin/sh\nprintf 'pipeline data\\n'\nexit 42\n")
		output, status := runCase(t, fakeBin, "produce | sh -c 'cat; exit 7'")
		if status != 7 || output != "pipeline data\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCommandNice(t *testing.T, adjustment, command string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration26MigrationName].Body.(string)
	if !ok || body != currentCommandNiceBody {
		t.Fatal("generated iteration-26 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1|-20,-15,-10,-5,0,5,10,15,19|}`, adjustment)
	return strings.ReplaceAll(body, `${2:command}`, command)
}

func TestCommandRenicePlaceholderContractAndPOSIXSyntax(t *testing.T) {
	if strings.Count(currentCommandReniceBody, `${1:processName}`) != 1 {
		t.Fatal("process-name placeholder changed")
	}
	if strings.Count(currentCommandReniceBody, `${2|-20,-15,-10,-5,0,5,10,15,19|}`) != 1 {
		t.Fatal("priority choice placeholder changed")
	}
	if !strings.Contains(currentCommandReniceBody, `pidof -- "${1:processName}"`) {
		t.Fatal("process name must remain one quoted pidof operand after option termination")
	}
	if strings.Contains(currentCommandReniceBody, "[[") || strings.Contains(currentCommandReniceBody, "function ") || strings.Contains(currentCommandReniceBody, "local ") {
		t.Fatal("renice snippet must remain POSIX sh")
	}
}

func TestCommandReniceBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, processName, pids, pidofStatus, failPID, failStatus, trailer string) (string, int) {
		t.Helper()
		fakeBin := t.TempDir()
		mustExecutable(t, filepath.Join(fakeBin, "pidof"), "#!/bin/sh\nprintf 'pidof-argc:<%s> name:<%s>\\n' \"$#\" \"$2\" >&2\n[ \"${SHELLMAN_PIDOF_STATUS:-0}\" -eq 0 ] || exit \"$SHELLMAN_PIDOF_STATUS\"\nprintf '%s\\n' \"$SHELLMAN_PIDS\"\n")
		mustExecutable(t, filepath.Join(fakeBin, "sudo"), "#!/bin/sh\nexec \"$@\"\n")
		mustExecutable(t, filepath.Join(fakeBin, "renice"), "#!/bin/sh\nprintf 'renice:<%s>:<%s>:<%s>:<%s>\\n' \"$1\" \"$2\" \"$3\" \"$4\"\nprintf 'renice-stderr:<%s>\\n' \"$4\" >&2\nif [ \"$4\" = \"$SHELLMAN_FAIL_PID\" ]; then exit \"$SHELLMAN_FAIL_STATUS\"; fi\n")
		body := runnableGeneratedCommandRenice(t, processName, "5")
		path := filepath.Join(t.TempDir(), "command-renice.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		environment := append(os.Environ(),
			"PATH="+fakeBin+":"+os.Getenv("PATH"),
			"SHELLMAN_PIDS="+pids,
			"SHELLMAN_PIDOF_STATUS="+pidofStatus,
			"SHELLMAN_FAIL_PID="+failPID,
			"SHELLMAN_FAIL_STATUS="+failStatus,
		)
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("all matching processes and streams", func(t *testing.T) {
		output, status := runCase(t, "", "worker", "101 202 303", "0", "", "0", `exit "$snippet_status"`)
		for _, want := range []string{"pidof-argc:<2> name:<worker>", "renice:<-n>:<5>:<-p>:<101>", "renice:<-n>:<5>:<-p>:<202>", "renice:<-n>:<5>:<-p>:<303>", "renice-stderr:<303>"} {
			if !strings.Contains(output, want) {
				t.Fatalf("output lacks %q: %q", want, output)
			}
		}
		if status != 0 {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, tc := range []struct {
		name   string
		status string
	}{
		{"no matching process", "1"},
		{"arbitrary pidof failure", "42"},
		{"pidof command-not-found status", "127"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output, status := runCase(t, "", "missing", "", tc.status, "", "0", `exit "$snippet_status"`)
			if status != mustAtoi(t, tc.status) || strings.Contains(output, "renice:") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("earlier renice failure is not masked and later PIDs are attempted", func(t *testing.T) {
		output, status := runCase(t, "", "worker", "101 202 303", "0", "202", "42", `exit "$snippet_status"`)
		if status != 42 || !strings.Contains(output, "renice:<-n>:<5>:<-p>:<303>") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("status one propagates", func(t *testing.T) {
		output, status := runCase(t, "", "worker", "101", "0", "101", "1", `exit "$snippet_status"`)
		if status != 1 || !strings.Contains(output, "renice:<-n>:<5>:<-p>:<101>") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("process name is one data operand", func(t *testing.T) {
		name := "-space * $HOME quote\\line\nsecond"
		input := `-space * \$HOME quote\\line
second`
		output, status := runCase(t, "", input, "808", "0", "", "0", `exit "$snippet_status"`)
		if status != 0 || !strings.Contains(output, "pidof-argc:<2> name:<"+name+">") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("caller variables and IFS are isolated", func(t *testing.T) {
		setup := `_shellman_renice_pids=caller; _shellman_renice_status=caller; _shellman_renice_pid=caller; IFS=:`
		trailer := `[ "$snippet_status" -eq 0 ] && [ "$_shellman_renice_pids" = caller ] && [ "$_shellman_renice_status" = caller ] && [ "$_shellman_renice_pid" = caller ] && [ "$IFS" = : ]`
		output, status := runCase(t, setup, "worker", "101 202", "0", "", "0", trailer)
		if status != 0 || !strings.Contains(output, "renice:<-n>:<5>:<-p>:<202>") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCommandRenice(t *testing.T, processName, priority string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration27MigrationName].Body.(string)
	if !ok || body != currentCommandReniceBody {
		t.Fatal("generated iteration-27 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `\$`, `$`)
	body = strings.ReplaceAll(body, `${1:processName}`, processName)
	return strings.ReplaceAll(body, `${2|-20,-15,-10,-5,0,5,10,15,19|}`, priority)
}

func mustAtoi(t *testing.T, value string) int {
	t.Helper()
	var result int
	if _, err := fmt.Sscanf(value, "%d", &result); err != nil {
		t.Fatal(err)
	}
	return result
}

func TestCommandRunPlaceholderContractAndPOSIXDocumentation(t *testing.T) {
	if strings.Count(currentCommandRunBody, `${1:result}`) != 1 {
		t.Fatal("result variable-name placeholder changed")
	}
	if strings.Count(currentCommandRunBody, `${2:command}`) != 1 {
		t.Fatal("raw command placeholder changed")
	}
	if !strings.HasPrefix(currentCommandRunBody, "# POSIX sh: command substitution runs in a subshell and removes trailing newlines.\n") {
		t.Fatal("POSIX comment must document the two non-obvious command-substitution semantics")
	}
	if !strings.HasSuffix(currentCommandRunBody, `${1:result}="$(${2:command})"`+"\n") {
		t.Fatal("historical quoted assignment or placeholder grammar changed")
	}
	if strings.Contains(currentCommandRunBody, "[[") || strings.Contains(currentCommandRunBody, "function ") || strings.Contains(currentCommandRunBody, "local ") {
		t.Fatal("command.run must remain POSIX sh")
	}
}

func TestCommandRunBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, destination, command, trailer string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandRun(t, destination, command)
		path := filepath.Join(t.TempDir(), "command-run.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("captures stdout preserves data and leaves stderr visible", func(t *testing.T) {
		command := `printf 'space * -n $HOME quote"backslash\\line\ninside\n\n'; printf 'visible stderr\n' >&2`
		trailer := `printf '<%s> status=<%s>\n' "$result" "$snippet_status"`
		output, status := runCase(t, "", "result", command, trailer)
		want := "visible stderr\n<space * -n $HOME quote\"backslash\\line\ninside> status=<0>\n"
		if status != 0 || output != want {
			t.Fatalf("status=%d output=%q want=%q", status, output, want)
		}
	})

	for _, tc := range []struct {
		name    string
		command string
		status  int
	}{
		{"status one", `printf captured; exit 1`, 1},
		{"arbitrary status", `printf captured; exit 42`, 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			trailer := `[ "$result" = captured ] && exit "$snippet_status"`
			output, status := runCase(t, "", "result", tc.command, trailer)
			if status != tc.status || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("command not found diagnostic and status propagate", func(t *testing.T) {
		output, status := runCase(t, "", "result", "shellman_command_that_does_not_exist", `exit "$snippet_status"`)
		if status != 127 || !strings.Contains(output, "shellman_command_that_does_not_exist") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("raw list uses final command status", func(t *testing.T) {
		command := `printf first; false; printf second; exit 42`
		trailer := `[ "$result" = firstsecond ] && exit "$snippet_status"`
		output, status := runCase(t, "", "result", command, trailer)
		if status != 42 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("pipeline uses POSIX last-command status", func(t *testing.T) {
		setup := `producer() { printf 'pipeline data\n'; return 42; }`
		command := `producer | sh -c 'cat; exit 7'`
		trailer := `[ "$result" = 'pipeline data' ] && exit "$snippet_status"`
		output, status := runCase(t, setup, "result", command, trailer)
		if status != 7 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("substitution side effects are isolated but assignment is in caller", func(t *testing.T) {
		setup := `marker=before; result=old; mutate() { marker=inside; printf changed; }`
		trailer := `[ "$snippet_status" -eq 0 ] && [ "$marker" = before ] && [ "$result" = changed ]`
		output, status := runCase(t, setup, "result", "mutate", trailer)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("caller IFS does not split or expand captured output", func(t *testing.T) {
		setup := `IFS=:`
		command := `printf 'one:two *'`
		trailer := `[ "$snippet_status" -eq 0 ] && [ "$IFS" = : ] && [ "$captured" = 'one:two *' ]`
		output, status := runCase(t, setup, "captured", command, trailer)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("empty command produces empty value", func(t *testing.T) {
		output, status := runCase(t, "", "result", "", `[ "$snippet_status" -eq 0 ] && [ -z "$result" ]`)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCommandRun(t *testing.T, destination, command string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration28MigrationName].Body.(string)
	if !ok || body != currentCommandRunBody {
		t.Fatal("generated iteration-28 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:result}`, destination)
	return strings.ReplaceAll(body, `${2:command}`, command)
}

func TestCommandSubstitutionPlaceholderContractAndPOSIXDocumentation(t *testing.T) {
	body := currentCommandSubstitutionBody
	if strings.Count(body, `${1:result}`) != 1 || strings.Count(body, `${2:command}`) != 1 {
		t.Fatal("command.substitution placeholders changed")
	}
	if !strings.HasPrefix(body, "# POSIX sh: command substitution runs in a subshell and removes trailing newlines.\n") {
		t.Fatal("POSIX comment must document subshell and trailing-newline semantics")
	}
	if !strings.HasSuffix(body, `${1:result}="$(${2:command})"`+"\n") {
		t.Fatal("quoted assignment or exact placeholder grammar changed")
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("command.substitution must remain POSIX sh")
	}
}

func TestCommandSubstitutionBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, destination, command, trailer string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandSubstitution(t, destination, command)
		path := filepath.Join(t.TempDir(), "command-substitution.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("capture transformations and stream scope", func(t *testing.T) {
		command := `printf 'space * -n $HOME quote"backslash\\line\ninternal\n\n'; printf 'stderr stays visible\n' >&2`
		trailer := `printf '<%s>:%s\n' "$captured" "$snippet_status"`
		output, status := runCase(t, "", "captured", command, trailer)
		want := "stderr stays visible\n<space * -n $HOME quote\"backslash\\line\ninternal>:0\n"
		if status != 0 || output != want {
			t.Fatalf("status=%d output=%q want=%q", status, output, want)
		}
	})

	for _, tc := range []struct {
		name    string
		command string
		status  int
	}{
		{"status one", `printf value; exit 1`, 1},
		{"arbitrary status", `printf value; exit 42`, 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output, status := runCase(t, "", "captured", tc.command, `[ "$captured" = value ] && exit "$snippet_status"`)
			if status != tc.status || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("command not found", func(t *testing.T) {
		output, status := runCase(t, "", "captured", "shellman_substitution_command_that_does_not_exist", `exit "$snippet_status"`)
		if status != 127 || !strings.Contains(output, "shellman_substitution_command_that_does_not_exist") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("raw list and final status", func(t *testing.T) {
		command := `printf first; false; printf second; exit 42`
		output, status := runCase(t, "", "captured", command, `[ "$captured" = firstsecond ] && exit "$snippet_status"`)
		if status != 42 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("pipeline status", func(t *testing.T) {
		setup := `producer() { printf 'pipeline output\n'; return 42; }`
		command := `producer | sh -c 'cat; exit 7'`
		output, status := runCase(t, setup, "captured", command, `[ "$captured" = 'pipeline output' ] && exit "$snippet_status"`)
		if status != 7 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("subshell isolation destination assignment and caller IFS", func(t *testing.T) {
		setup := `marker=caller; captured=old; IFS=:`
		command := `marker=substitution; printf 'one:two *'`
		trailer := `[ "$snippet_status" -eq 0 ] && [ "$marker" = caller ] && [ "$captured" = 'one:two *' ] && [ "$IFS" = : ]`
		output, status := runCase(t, setup, "captured", command, trailer)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("empty command", func(t *testing.T) {
		output, status := runCase(t, "", "captured", "", `[ "$snippet_status" -eq 0 ] && [ -z "$captured" ]`)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCommandSubstitution(t *testing.T, destination, command string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration29MigrationName].Body.(string)
	if !ok || body != currentCommandSubstitutionBody {
		t.Fatal("generated iteration-29 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:result}`, destination)
	return strings.ReplaceAll(body, `${2:command}`, command)
}

func TestCommandSuccessCheckPlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := strings.Join(anyStrings(t, currentCommandSuccessCheckBody), "\n")
	if strings.Count(body, `${1:command}`) != 1 {
		t.Fatal("raw command placeholder changed")
	}
	if !strings.Contains(body, "if {\n\t${1:command}\n} >/dev/null 2>&1; then") {
		t.Fatal("redirection must cover the complete raw command placeholder")
	}
	if !strings.Contains(body, `set -- "\$?"`) || !strings.Contains(body, `exit "\$1"`) {
		t.Fatal("failure status must be captured and returned inside isolated positional parameters")
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("success-check must remain POSIX sh")
	}
}

func TestCommandSuccessCheckBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, command, trailer string) (string, int) {
		t.Helper()
		body := runnableGeneratedCommandSuccessCheck(t, command)
		path := filepath.Join(t.TempDir(), "command-success-check.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("success reports success suppresses checked streams and returns zero", func(t *testing.T) {
		command := `printf 'first stdout\n'; printf 'first stderr\n' >&2; printf 'last stdout\n'; printf 'last stderr\n' >&2`
		output, status := runCase(t, "", command, `exit "$snippet_status"`)
		if status != 0 || output != "succeed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, tc := range []struct {
		name    string
		command string
		status  int
	}{
		{"status one", "false", 1},
		{"arbitrary status", `fail() { return 42; }; fail`, 42},
	} {
		t.Run(tc.name, func(t *testing.T) {
			output, status := runCase(t, "", tc.command, `exit "$snippet_status"`)
			if status != tc.status || output != "failed\n" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("command not found is quiet and preserves 127", func(t *testing.T) {
		output, status := runCase(t, "", "shellman_success_command_that_does_not_exist", `exit "$snippet_status"`)
		if status != 127 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("raw list final status and complete redirection scope", func(t *testing.T) {
		command := `printf 'hidden one\n'; printf 'hidden two\n' >&2; sh -c 'exit 42'`
		output, status := runCase(t, "", command, `exit "$snippet_status"`)
		if status != 42 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("pipeline uses POSIX last-command status and group redirection", func(t *testing.T) {
		command := `printf 'hidden producer stderr\n' >&2 | sh -c 'cat; exit 7'`
		output, status := runCase(t, "", command, `exit "$snippet_status"`)
		if status != 7 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("checked side effects remain in caller and positional parameters are preserved", func(t *testing.T) {
		setup := `set -- caller one; marker=before; check() { marker='after value'; return 42; }`
		trailer := `[ "$snippet_status" -eq 42 ] && [ "$marker" = 'after value' ] && [ "$#" -eq 2 ] && [ "$1" = caller ] && [ "$2" = one ]`
		output, status := runCase(t, setup, "check", trailer)
		if status != 0 || output != "failed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("quoted arguments retain data under caller IFS", func(t *testing.T) {
		setup := `IFS=:; check() { [ "$1" = 'space * -n $HOME quote"backslash\line' ]; }`
		command := `check 'space * -n $HOME quote"backslash\line'`
		trailer := `[ "$snippet_status" -eq 0 ] && [ "$IFS" = : ]`
		output, status := runCase(t, setup, command, trailer)
		if status != 0 || output != "succeed\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("diagnostic output failure is not masked", func(t *testing.T) {
		body := runnableGeneratedCommandSuccessCheck(t, "false")
		path := filepath.Join(t.TempDir(), "command-success-check-output.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+body))
		_, err := runCommand(".", nil, "sh", "-c", `sh "$1" >/dev/full`, "sh", path)
		if exitCode(err) == 0 {
			t.Fatal("failed diagnostic write was masked")
		}
	})
}

func runnableGeneratedCommandSuccessCheck(t *testing.T, command string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration30MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentCommandSuccessCheckBody) {
		t.Fatal("generated iteration-30 body differs from candidate")
	}
	result := strings.Join(anyStrings(t, body), "\n")
	result = strings.ReplaceAll(result, `\$`, `$`)
	return strings.ReplaceAll(result, `${1:command}`, command)
}

func TestCryptographyBase64DecodePlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := currentCryptographyBase64DecodeBody
	for _, placeholder := range []string{`${1:base64Decoded}`, `${2|stringToDecode,${variableToDecode}|}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q count = %d, want 1", placeholder, strings.Count(body, placeholder))
		}
	}
	if !strings.Contains(body, `printf '%s' "${2|stringToDecode,${variableToDecode}|}" | base64 -d`) {
		t.Fatal("encoded data must be printed exactly and piped to the decoder")
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("base64 decoder snippet must remain POSIX sh syntax")
	}
}

func TestCryptographyBase64DecodeBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "base64")

	runCase := func(t *testing.T, setup, destination, encoded, trailer string, environment []string) (string, int) {
		t.Helper()
		body := runnableGeneratedCryptographyBase64Decode(t, destination, encoded)
		path := filepath.Join(t.TempDir(), "base64-decode.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	tests := []struct {
		name string
		data []byte
	}{
		{"ordinary text", []byte("Shellman base64\n")},
		{"whitespace glob hyphen quote backslash and expansion-looking text", []byte("space * -n 'quote' backslash\\ $HOME\ninternal newline\n")},
		{"leading newline", []byte("\nleading\n")},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			encoded := base64.StdEncoding.EncodeToString(tc.data)
			trailer := `[ "$snippet_status" -eq 0 ] && [ "$decoded" = "$expected" ]`
			environment := append(os.Environ(), "expected="+strings.TrimRight(string(tc.data), "\n"))
			output, status := runCase(t, ``, "decoded", `"${variableToDecode}"`, trailer, append(environment, "variableToDecode="+encoded))
			if status != 0 || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("literal nested choice selection", func(t *testing.T) {
		encoded := base64.StdEncoding.EncodeToString([]byte("literal choice"))
		output, status := runCase(t, "", "decoded", `"`+encoded+`"`, `[ "$snippet_status" -eq 0 ] && [ "$decoded" = 'literal choice' ]`, nil)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("invalid input propagates decoder status and diagnostic", func(t *testing.T) {
		output, status := runCase(t, "", "decoded", `"%%%"`, `exit "$snippet_status"`, nil)
		if status != 1 || !strings.Contains(output, "invalid") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("arbitrary decoder failure propagates", func(t *testing.T) {
		fakeBin := t.TempDir()
		mustExecutable(t, filepath.Join(fakeBin, "base64"), "#!/bin/sh\ncat >/dev/null\nprintf partial\nprintf 'decoder error\\n' >&2\nexit 42\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		output, status := runCase(t, "", "decoded", `"QQ=="`, `[ "$decoded" = partial ] && exit "$snippet_status"`, environment)
		if status != 42 || !strings.Contains(output, "decoder error") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("decoder command not found propagates", func(t *testing.T) {
		fakeBin := t.TempDir()
		output, status := runCase(t, "", "decoded", `"QQ=="`, `exit "$snippet_status"`, []string{"PATH=" + fakeBin})
		if status != 127 || !strings.Contains(output, "base64") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("empty value decodes to empty output", func(t *testing.T) {
		output, status := runCase(t, "", "decoded", `""`, `[ "$snippet_status" -eq 0 ] && [ -z "$decoded" ]`, nil)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCryptographyBase64Decode(t *testing.T, destination, encodedExpression string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration31MigrationName].Body.(string)
	if !ok || body != currentCryptographyBase64DecodeBody {
		t.Fatal("generated iteration-31 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `\$`, `$`)
	body = strings.ReplaceAll(body, `${1:base64Decoded}`, destination)
	return strings.ReplaceAll(body, `${2|stringToDecode,${variableToDecode}|}`, encodedExpression)
}

func TestCryptographyBase64EncodePlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := currentCryptographyBase64EncodeBody
	for _, placeholder := range []string{`${1:base64Encoded}`, `${2|stringToEncode,${variableToEncode}|}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q count = %d, want 1", placeholder, strings.Count(body, placeholder))
		}
	}
	if !strings.Contains(body, `printf '%s' "${2|stringToEncode,${variableToEncode}|}" | base64`) {
		t.Fatal("input data must be printed exactly and piped to the encoder")
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("base64 encoder snippet must remain POSIX sh syntax")
	}
}

func TestCryptographyBase64EncodeBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "base64")

	runCase := func(t *testing.T, setup, destination, inputExpression, trailer string, environment []string) (string, int) {
		t.Helper()
		body := runnableGeneratedCryptographyBase64Encode(t, destination, inputExpression)
		path := filepath.Join(t.TempDir(), "base64-encode.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	tests := []struct {
		name string
		data string
	}{
		{"ordinary text", "Shellman base64"},
		{"empty input", ""},
		{"whitespace glob hyphen quote backslash and expansion-looking text", " space * -n 'quote' backslash\\ $HOME "},
		{"leading internal and trailing newlines", "\nleading\ninternal\n"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			want := base64.StdEncoding.EncodeToString([]byte(tc.data))
			trailer := `[ "$snippet_status" -eq 0 ] && [ "$encoded" = "$expected" ]`
			environment := append(os.Environ(), "variableToEncode="+tc.data, "expected="+want, "IFS=*,")
			output, status := runCase(t, "", "encoded", `${variableToEncode}`, trailer, environment)
			if status != 0 || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("literal choice selection", func(t *testing.T) {
		want := base64.StdEncoding.EncodeToString([]byte("stringToEncode"))
		output, status := runCase(t, "", "encoded", `stringToEncode`, `[ "$snippet_status" -eq 0 ] && [ "$encoded" = "$expected" ]`, []string{"expected=" + want})
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, failureStatus := range []int{1, 42} {
		t.Run(fmt.Sprintf("encoder failure %d preserves partial output status and stderr", failureStatus), func(t *testing.T) {
			fakeBin := t.TempDir()
			mustExecutable(t, filepath.Join(fakeBin, "base64"), fmt.Sprintf("#!/bin/sh\ncat >/dev/null\nprintf 'partial\\n\\n'\nprintf 'encoder error\\n' >&2\nexit %d\n", failureStatus))
			environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
			output, status := runCase(t, "", "encoded", `data`, `[ "$encoded" = partial ] && exit "$snippet_status"`, environment)
			if status != failureStatus || !strings.Contains(output, "encoder error") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("encoder command not found propagates", func(t *testing.T) {
		fakeBin := t.TempDir()
		output, status := runCase(t, "", "encoded", `data`, `exit "$snippet_status"`, []string{"PATH=" + fakeBin})
		if status != 127 || !strings.Contains(output, "base64") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("destination assignment failure is nonzero", func(t *testing.T) {
		output, status := runCase(t, "readonly encoded", "encoded", `data`, `exit "$snippet_status"`, nil)
		if status == 0 || !strings.Contains(output, "encoded") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCryptographyBase64Encode(t *testing.T, destination, inputExpression string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration32MigrationName].Body.(string)
	if !ok || body != currentCryptographyBase64EncodeBody {
		t.Fatal("generated iteration-32 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `\$`, `$`)
	body = strings.ReplaceAll(body, `${1:base64Encoded}`, destination)
	return strings.ReplaceAll(body, `${2|stringToEncode,${variableToEncode}|}`, inputExpression)
}

func TestCryptographyHashPlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := currentCryptographyHashBody
	for _, placeholder := range []string{`${1:hash}`, `${2:variableToHash}`, `${3|md5sum,shasum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum|}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q count = %d, want 1", placeholder, strings.Count(body, placeholder))
		}
	}
	if !strings.Contains(body, `printf '%s' "\$${2:variableToHash}" | ${3|md5sum,shasum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum|}`) {
		t.Fatal("variable contents must be produced exactly and sent to the selected hash utility")
	}
	if strings.Contains(body, "[[") || strings.Contains(body, "function ") || strings.Contains(body, "local ") {
		t.Fatal("hash snippet must remain POSIX sh syntax")
	}
}

func TestCryptographyHashBehavior(t *testing.T) {
	requireCommand(t, "sh")

	runCase := func(t *testing.T, setup, destination, variable, utility, trailer string, environment []string) (string, int) {
		t.Helper()
		body := runnableGeneratedCryptographyHash(t, destination, variable, utility)
		path := filepath.Join(t.TempDir(), "hash.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	utilities := []struct {
		name string
		hash func([]byte) []byte
	}{
		{"md5sum", func(data []byte) []byte { sum := md5.Sum(data); return sum[:] }},
		{"shasum", func(data []byte) []byte { sum := sha1.Sum(data); return sum[:] }},
		{"sha1sum", func(data []byte) []byte { sum := sha1.Sum(data); return sum[:] }},
		{"sha224sum", func(data []byte) []byte { sum := sha256.Sum224(data); return sum[:] }},
		{"sha256sum", func(data []byte) []byte { sum := sha256.Sum256(data); return sum[:] }},
		{"sha384sum", func(data []byte) []byte { sum := sha512.Sum384(data); return sum[:] }},
		{"sha512sum", func(data []byte) []byte { sum := sha512.Sum512(data); return sum[:] }},
	}
	inputs := []struct {
		name string
		data string
	}{
		{"ordinary text", "Shellman hash"},
		{"empty input", ""},
		{"whitespace glob hyphen quote backslash and expansion-looking text", " -n * 'quoted' \\ $HOME "},
		{"leading internal and trailing newlines", "\nleading\ninternal\n"},
	}
	for _, utility := range utilities {
		if _, err := exec.LookPath(utility.name); err != nil {
			t.Logf("SKIP utility %s: not installed", utility.name)
			continue
		}
		for _, input := range inputs {
			t.Run(utility.name+"/"+input.name, func(t *testing.T) {
				want := hex.EncodeToString(utility.hash([]byte(input.data)))
				trailer := `[ "$snippet_status" -eq 0 ] && [ "$hash_result" = "$expected" ]`
				environment := append(os.Environ(), "variableToHash="+input.data, "expected="+want, "IFS=*,")
				output, status := runCase(t, "", "hash_result", "variableToHash", utility.name, trailer, environment)
				if status != 0 || output != "" {
					t.Fatalf("status=%d output=%q", status, output)
				}
			})
		}
	}

	t.Run("command choice is substituted literally", func(t *testing.T) {
		fakeBin := t.TempDir()
		mustExecutable(t, filepath.Join(fakeBin, "hash_choice"), "#!/bin/sh\ncat >/dev/null\nprintf 'choice-digest  -\\n'\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		output, status := runCase(t, "", "hash_result", "variableToHash", "hash_choice", `[ "$snippet_status" -eq 0 ] && [ "$hash_result" = choice-digest ]`, environment)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, failureStatus := range []int{1, 42} {
		t.Run(fmt.Sprintf("utility failure %d propagates and discards partial result", failureStatus), func(t *testing.T) {
			fakeBin := t.TempDir()
			mustExecutable(t, filepath.Join(fakeBin, "hash_failure"), fmt.Sprintf("#!/bin/sh\ncat >/dev/null\nprintf 'partial-digest  -\\n'\nprintf 'hash error\\n' >&2\nexit %d\n", failureStatus))
			environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
			output, status := runCase(t, "", "hash_result", "variableToHash", "hash_failure", `[ -z "$hash_result" ] && exit "$snippet_status"`, environment)
			if status != failureStatus || !strings.Contains(output, "hash error") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("utility command not found propagates 127", func(t *testing.T) {
		fakeBin := t.TempDir()
		output, status := runCase(t, "", "hash_result", "variableToHash", "hash_missing", `exit "$snippet_status"`, []string{"PATH=" + fakeBin})
		if status != 127 || !strings.Contains(output, "hash_missing") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("destination assignment failure is nonzero", func(t *testing.T) {
		output, status := runCase(t, "readonly hash_result", "hash_result", "variableToHash", "sha256sum", `exit "$snippet_status"`, nil)
		if status == 0 || !strings.Contains(output, "hash_result") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedCryptographyHash(t *testing.T, destination, variable, utility string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration33MigrationName].Body.(string)
	if !ok || body != currentCryptographyHashBody {
		t.Fatal("generated iteration-33 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `\$`, `$`)
	body = strings.ReplaceAll(body, `${1:hash}`, destination)
	body = strings.ReplaceAll(body, `${2:variableToHash}`, variable)
	return strings.ReplaceAll(body, `${3|md5sum,shasum,sha1sum,sha224sum,sha256sum,sha384sum,sha512sum|}`, utility)
}

func TestDateNowShortPlaceholderContractAndPOSIXSyntax(t *testing.T) {
	body := currentDateNowShortBody
	for _, placeholder := range []string{`${1:dateShort}`, `${0:# format: yyyy/mm/dd}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q count = %d, want 1", placeholder, strings.Count(body, placeholder))
		}
	}
	if !strings.Contains(body, `date '+%Y/%m/%d'`) {
		t.Fatal("date must use the documented year/month/day format")
	}
	if strings.Contains(body, "date -I") {
		t.Fatal("GNU-specific date -I must not remain")
	}
}

func TestDateNowShortBehaviorAndStatus(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "date")

	runCase := func(t *testing.T, setup, destination, trailer string, environment []string) (string, int) {
		t.Helper()
		body := runnableGeneratedDateNowShort(t, destination, "# format: yyyy/mm/dd")
		path := filepath.Join(t.TempDir(), "date-now-short.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+body+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", environment, "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("installed date emits documented slash format", func(t *testing.T) {
		output, status := runCase(t, "", "date_result", `
case $date_result in
  [0-9][0-9][0-9][0-9]/[0-9][0-9]/[0-9][0-9]) [ "$snippet_status" -eq 0 ] ;;
  *) exit 1 ;;
esac`, append(os.Environ(), "TZ=UTC0", "LC_ALL=C"))
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("format operand and final-tabstop comment expand exactly", func(t *testing.T) {
		fakeBin := t.TempDir()
		mustExecutable(t, filepath.Join(fakeBin, "date"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = '+%Y/%m/%d' ] || exit 42\nprintf '2034/05/06\\n'\n")
		environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
		output, status := runCase(t, "", "date_result", `[ "$snippet_status" -eq 0 ] && [ "$date_result" = 2034/05/06 ]`, environment)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, failureStatus := range []int{1, 42} {
		t.Run(fmt.Sprintf("date failure %d propagates status and exposes partial assignment", failureStatus), func(t *testing.T) {
			fakeBin := t.TempDir()
			mustExecutable(t, filepath.Join(fakeBin, "date"), fmt.Sprintf("#!/bin/sh\nprintf 'partial-date\\n'\nprintf 'date error\\n' >&2\nexit %d\n", failureStatus))
			environment := append(os.Environ(), "PATH="+fakeBin+":"+os.Getenv("PATH"))
			output, status := runCase(t, "", "date_result", `[ "$date_result" = partial-date ] && exit "$snippet_status"`, environment)
			if status != failureStatus || !strings.Contains(output, "date error") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("date command not found propagates 127", func(t *testing.T) {
		fakeBin := t.TempDir()
		output, status := runCase(t, "", "date_result", `exit "$snippet_status"`, []string{"PATH=" + fakeBin})
		if status != 127 || !strings.Contains(output, "date") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("destination assignment failure is nonzero", func(t *testing.T) {
		output, status := runCase(t, "readonly date_result", "date_result", `exit "$snippet_status"`, nil)
		if status == 0 || !strings.Contains(output, "date_result") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func runnableGeneratedDateNowShort(t *testing.T, destination, finalComment string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration34MigrationName].Body.(string)
	if !ok || body != currentDateNowShortBody {
		t.Fatal("generated iteration-34 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `\$`, `$`)
	body = strings.ReplaceAll(body, `${1:dateShort}`, destination)
	return strings.ReplaceAll(body, `${0:# format: yyyy/mm/dd}`, finalComment)
}

func TestArraySetElementAtPlaceholderContractAndBashDocumentation(t *testing.T) {
	body := currentArraySetElementAtBody
	for _, placeholder := range []string{`${1:myArray}`, `${2:index}`, `${3:value}`} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q count = %d, want 1", placeholder, strings.Count(body, placeholder))
		}
	}
	if !strings.HasPrefix(body, "# Bash-only: indexed arrays and arithmetic array indices are not specified by POSIX sh.\n") {
		t.Fatal("Bash-only comment must precisely identify indexed-array and arithmetic-index dependencies")
	}
	if !strings.HasSuffix(body, `${1:myArray}[${2:index}]="${3:value}"`+"\n") {
		t.Fatal("historical quoted assignment and exact placeholders changed")
	}
}

func TestArraySetElementAtBehavior(t *testing.T) {
	requireCommand(t, "bash")

	runCase := func(t *testing.T, setup, index, value, assertion string) (string, int) {
		t.Helper()
		body := runnableGeneratedArraySetElementAt(t, "array", index, value)
		path := filepath.Join(t.TempDir(), "array-set-element-at.sh")
		mustWriteFile(t, path, []byte("#!/usr/bin/env bash\n"+setup+"\n"+body+"snippet_status=$?\n"+assertion+"\n"))
		run(t, ".", "bash", "-n", path)
		output, err := runCommand(".", nil, "bash", path)
		return string(output), exitCode(err)
	}

	tests := []struct {
		name      string
		setup     string
		index     string
		value     string
		assertion string
	}{
		{"dense existing index is updated", `array=(zero one two)`, "1", "changed", `[[ $snippet_status -eq 0 && ${#array[@]} -eq 3 && ${array[0]} == zero && ${array[1]} == changed && ${array[2]} == two ]]`},
		{"sparse existing index is updated", `array=(); array[2]=two; array[8]=eight`, "8", "changed", `[[ $snippet_status -eq 0 && ${#array[@]} -eq 2 && ${array[2]} == two && ${array[8]} == changed && ! -v 'array[0]' ]]`},
		{"beyond highest index creates sparse element", `array=(zero one)`, "20", "far", `[[ $snippet_status -eq 0 && ${#array[@]} -eq 3 && ${array[20]} == far && ! -v 'array[2]' ]]`},
		{"index zero on empty array", `array=()`, "0", "first", `[[ $snippet_status -eq 0 && ${#array[@]} -eq 1 && ${array[0]} == first ]]`},
		{"previously unset interior index", `array=(); array[0]=zero; array[3]=three`, "2", "two", `[[ $snippet_status -eq 0 && ${#array[@]} -eq 3 && ${array[2]} == two && ! -v 'array[1]' ]]`},
		{"empty replacement", `array=(old)`, "0", "", `[[ $snippet_status -eq 0 && -v 'array[0]' && -z ${array[0]} ]]`},
		{"special characters preserve one value", `array=(old)`, "0", `space * -n \"quote\" a\\b
\$HOME`, `[[ $snippet_status -eq 0 && ${#array[@]} -eq 1 && ${array[0]} == $'space * -n "quote" a\\b\n$HOME' ]]`},
		{"arithmetic expression index", `array=(zero one two three)`, "1+2", "changed", `[[ $snippet_status -eq 0 && ${array[3]} == changed && ${#array[@]} -eq 4 ]]`},
		{"negative index updates last element", `array=(zero one two)`, "-1", "last", `[[ $snippet_status -eq 0 && ${#array[@]} -eq 3 && ${array[2]} == last ]]`},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output, status := runCase(t, test.setup, test.index, test.value, test.assertion)
			if status != 0 {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("negative index on empty array fails", func(t *testing.T) {
		output, status := runCase(t, `array=()`, "-1", "value", `exit "$snippet_status"`)
		if status == 0 || !strings.Contains(output, "bad array subscript") {
			t.Fatalf("failure was masked: status=%d output=%q", status, output)
		}
	})

	t.Run("readonly array assignment failure propagates", func(t *testing.T) {
		output, status := runCase(t, `array=(old); readonly -a array`, "0", "new", `exit "$snippet_status"`)
		if status == 0 || !strings.Contains(output, "readonly variable") {
			t.Fatalf("failure was masked: status=%d output=%q", status, output)
		}
	})
}

func TestArrayNamespaceIterationInventory(t *testing.T) {
	ordered, err := readSnippets(rootDirectory)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"array.all-elements", "array.at-index", "array.concat", "array.contains",
		"array.declare", "array.delete-at", "array.delete", "array.filter",
		"array.iterate", "array.length", "array.print", "array.push",
		"array.range", "array.replace", "array.reverse", "array.set-element-at",
	}
	got := make([]string, 0, len(want))
	for _, item := range ordered {
		if item.namespace == "array" {
			got = append(got, item.name)
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("array traversal = %#v, want %#v", got, want)
	}
	if got := ordered[22].name; got != "command.failure-check" {
		t.Fatalf("next traversal entry = %q, want command.failure-check", got)
	}
}

func runnableGeneratedArraySetElementAt(t *testing.T, array, index, value string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration22MigrationName].Body.(string)
	if !ok || body != currentArraySetElementAtBody {
		t.Fatal("generated iteration-22 body differs from candidate")
	}
	body = strings.ReplaceAll(body, `${1:myArray}`, array)
	body = strings.ReplaceAll(body, `${2:index}`, index)
	return strings.ReplaceAll(body, `${3:value}`, value)
}

func runWithEnvironment(t *testing.T, dir string, environment []string, name string, args ...string) string {
	t.Helper()
	output, err := runCommand(dir, environment, name, args...)
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return string(output)
}

func assertDirectoryEmpty(t *testing.T, path string) {
	t.Helper()
	entries, err := os.ReadDir(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("directory %s is not empty: %v", path, entries)
	}
}

func requireCommand(t *testing.T, name string) {
	t.Helper()
	if _, err := exec.LookPath(name); err != nil {
		t.Fatalf("required command %q is unavailable: %v", name, err)
	}
}

func scriptPath(t *testing.T, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "compress-tar-gz.sh")
	if err := os.WriteFile(path, []byte("#!/bin/sh\nset -eu\n"+body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func run(t *testing.T, dir, name string, args ...string) string {
	t.Helper()
	output, err := runCommand(dir, nil, name, args...)
	if err != nil {
		t.Fatalf("%s %s failed: %v\n%s", name, strings.Join(args, " "), err, output)
	}
	return string(output)
}

func runCommand(dir string, environment []string, name string, args ...string) ([]byte, error) {
	command := exec.Command(name, args...)
	command.Dir = dir
	if environment != nil {
		command.Env = environment
	}
	return command.CombinedOutput()
}

func exitCode(err error) int {
	if err == nil {
		return 0
	}
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}
	return -1
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatal(err)
	}
}

func mustExecutable(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o755); err != nil {
		t.Fatal(err)
	}
}

func TestV6SnippetInventoryAndOrder(t *testing.T) {
	ordered, err := readSnippets(rootDirectory)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(ordered), 278; got != want {
		t.Fatalf("snippet count = %d, want %d", got, want)
	}

	var namespaces []string
	for _, item := range ordered {
		if len(namespaces) == 0 || namespaces[len(namespaces)-1] != item.namespace {
			namespaces = append(namespaces, item.namespace)
		}
	}
	if !reflect.DeepEqual(namespaces, v6Namespaces) {
		t.Fatalf("namespaces = %q, want %q", namespaces, v6Namespaces)
	}

	names := sortedSnippetNames(ordered)
	if !sort.StringsAreSorted(names) {
		t.Fatal("generated snippet names are not sorted")
	}

	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	if got := jsonObjectKeys(t, generated); !reflect.DeepEqual(got, names) {
		t.Fatal("JSON key order differs from the v6 snippet-name order")
	}
}

func TestGenerateIsDeterministic(t *testing.T) {
	firstSnippets, firstCommands, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	secondSnippets, secondCommands, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(firstSnippets, secondSnippets) || !bytes.Equal(firstCommands, secondCommands) {
		t.Fatal("repeated generation produced different output")
	}
}

func assertFileBytes(t *testing.T, path string, generated []byte) {
	t.Helper()
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(generated, want) {
		t.Fatalf("generated output differs from %s", path)
	}
}

func assertSHA256(t *testing.T, name string, data []byte, want string) {
	t.Helper()
	digest := sha256.Sum256(data)
	if got := hex.EncodeToString(digest[:]); got != want {
		t.Fatalf("%s SHA-256 = %s, want approved contract %s", name, got, want)
	}
}

func jsonObjectKeys(t *testing.T, data []byte) []string {
	t.Helper()
	decoder := json.NewDecoder(bytes.NewReader(data))
	token, err := decoder.Token()
	if err != nil || token != json.Delim('{') {
		t.Fatalf("decode generated JSON object: %v", err)
	}

	var keys []string
	for decoder.More() {
		token, err := decoder.Token()
		if err != nil {
			t.Fatal(err)
		}
		keys = append(keys, token.(string))
		var value json.RawMessage
		if err := decoder.Decode(&value); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := decoder.Token(); err != nil {
		t.Fatal(err)
	}
	return keys
}

func TestGeneratedFileStructure(t *testing.T) {
	if snippetOutputPath != "snippets/snippets.json" {
		t.Fatalf("unexpected snippet output path %q", snippetOutputPath)
	}
	if documentOutputPath != "COMMANDS.md" {
		t.Fatalf("unexpected documentation output path %q", documentOutputPath)
	}
	if strings.ContainsAny(snippetOutputPath, "\\") {
		t.Fatalf("snippet output path is not portable: %q", snippetOutputPath)
	}
}
