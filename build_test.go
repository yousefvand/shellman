package main

import (
	"bytes"
	"context"
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
	"slices"
	"sort"
	"strings"
	"testing"
	"time"
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
	currentV7Iteration35SHA256          = "f4a6f9ea4308ff3ee51e8eb360b222352dbae37843281d2c62d05d6e48665225"
	currentV7Iteration36SHA256          = "cc808323680c56c84ae918b9c40c631fe63a3f2e518b4fdafc21f353e1984d78"
	currentV7Iteration37SHA256          = "a5f648c681c9a081e5b99ae087dd1969a9e26b17ef8cf55ccc3da8d213b06ae2"
	currentV7Iteration38SHA256          = "b9bbcc604a86e8d314ed4d7f0ee235491667b6634084d6a768524ece5e1a79fe"
	currentV7Iteration39SHA256          = "d34ab5cb9e89e065b2af8213be8f83c648966acd9482d1492ae662999890f1c9"
	currentV7Iteration40SHA256          = "59aba95d557ad0490c8effae48a5531022d443d420c0ed72c30b05afd97d3e8a"
	currentV7Iteration41SHA256          = "583a23e13513bbf052587c5576103bd3d1bc636b24665501bcf7038f92221f1b"
	currentV7Iteration42SHA256          = "ff586d0e36a95bc676441b850d1f3edc1fc8d113b68e308b69965c7be08316d4"
	currentV7Iteration43SHA256          = "8cb8126c68afe47860b489dc24ea623dd9ba440844ff9df4a51b3aa90af22554"
	currentV7Iteration44SHA256          = "1be38e1f47b9fc3a978de6af6df463daa949a4656c6e91ddee7bb281575943f9"
	currentV7Iteration45SHA256          = "4dddc3b8e76a28162d39590fd78f4d5d2790a5feba49971da77951baa79378e6"
	currentV7Iteration46SHA256          = "d33608de062dcc979cdec9721608e377b18f270ca4ef96f2388b8b3eb9675587"
	currentV7Iteration47SHA256          = "39b208febbdbc21217d308393d172b3660b29061fb208fb5a39be60e013486a1"
	currentV7Iteration48SHA256          = "62426db99ad2860ad55049ee9c2fb23691dacc44def107ac5b0615228e3b505b"
	currentV7Iteration49SHA256          = "c91a7b5ec85ccc87b1840be2dc58d6da1d17c071687bec086708c82f6cabe2f8"
	currentV7Iteration50SHA256          = "1ce3b0590aa0cfa6b46073461e0b987620a08b6a598e26fb91bcb380d753161b"
	currentV7Iteration51SHA256          = "c4ae28b18fa11e795a1be93f2962f0612cdeb04ebb44cf53d8c8102ed33238e7"
	currentV7Iteration52SHA256          = "cff68abcc06f3f9a3b9698926e8b76d29a5f1eaec11da70cb82fc1afd25335e2"
	currentV7Iteration53SHA256          = "65a611dcf0e3f9bc4f98940272f4168b4f382b9b019ab43b9392aa38f26cc265"
	currentV7Iteration54SHA256          = "8d9e80dffd00ee44106d8f0ffc3ba2bce8f9f24ce07cdac69a28b14fd3d12f07"
	currentV7Iteration55SHA256          = "c9f1dae54f2d911218fab6d88df4ffd143c05336b19e195b87146fa7f28a86e2"
	currentV7Iteration56SHA256          = "51f23ce80d31b3f020e1c6ecfa32bec4bb1ffa8bfc17492d6cb9407663a54db4"
	currentV7Iteration58SHA256          = "e3b820a02fe36079e8caff51dcca8828513a5408439432c4ef56ead7c73e81b9"
	currentV7Iteration66SHA256          = "d0b12211fbc9febf6ec9c56728f2f6b55d6c1cb8e64a90812f2694c8071c2a0b"
	currentV7Iteration68SHA256          = "d7b86fd84d520ae1021e1a024ee90cfe56b1047d11a83f01639613a4bc4b4c72"
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
	iteration35MigrationName            = "date.day-of-month-current"
	iteration36MigrationName            = "date.day-of-week-current"
	iteration37MigrationName            = "date.day-of-year-current"
	iteration38MigrationName            = "date.local-short"
	iteration39MigrationName            = "date.month-name-current"
	iteration40MigrationName            = "date.month-number-current"
	iteration41MigrationName            = "date.utc-long"
	iteration42MigrationName            = "date.year-current"
	iteration43MigrationName            = "event.on-ctrl-c"
	iteration44MigrationName            = "event.on-exit"
	iteration45MigrationName            = "filesystem.directories-iterate"
	iteration46MigrationName            = "filesystem.directory-create-nested"
	iteration47MigrationName            = "filesystem.directory-create"
	iteration48MigrationName            = "filesystem.directory-delete-nested"
	iteration49MigrationName            = "filesystem.file-delete"
	iteration50MigrationName            = "filesystem.file-read"
	v6FileDeleteBody                    = "rm -f ${1|/path/to/file,${pathToFile}|}\n"
	currentFileDeleteBody               = "file_path=\"${1|/path/to/file,${pathToFile}|}\"\ncase $file_path in\n  -*) file_path=./$file_path ;;\nesac\nrm -f \"$file_path\"\n"
	v6DirectoryDeleteNestedBody         = "rm -rf ${1|/path/to/directory,${pathToDirectory}|}\n"
	currentDirectoryDeleteNestedBody    = "directory_path=\"${1|/path/to/directory,${pathToDirectory}|}\"\ncase $directory_path in\n  -*) directory_path=./$directory_path ;;\nesac\nrm -rf \"$directory_path\"\n"
	v6DirectoryCreateBody               = "mkdir \"${1:directory name}\"\n"
	currentDirectoryCreateBody          = "directory_path=\"${1:directory name}\"\ncase $directory_path in\n  -*) directory_path=./$directory_path ;;\nesac\nmkdir \"$directory_path\"\n"
	v6DirectoryCreateNestedBody         = "mkdir -p ${1|\"parent dir/child dir\",\"${pathToCreate}\"|}\n"
	currentDirectoryCreateNestedBody    = "directory_path=${1|\"parent dir/child dir\",\"${pathToCreate}\"|}\ncase $directory_path in\n  -*) directory_path=./$directory_path ;;\nesac\nmkdir -p \"$directory_path\"\n"
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
	v6DayOfMonthCurrentBody             = "${1:dayOfMonth}=\\$(date +%d)\n"
	currentDayOfMonthCurrentBody        = "# date %d renders the 1..31 day of month as 01..31.\n${1:dayOfMonth}=\\$(date +%d)\n"
	v6DayOfWeekCurrentBody              = "${1:dayOfWeek}=\\$(date +%${2|A,a|})\n"
	currentDayOfWeekCurrentBody         = "# Weekday names follow the current locale.\n${1:dayOfWeek}=\\$(date +%${2|A,a|})\n"
	v6DayOfYearCurrentBody              = "${1:dayOfYear}=\\$(date +%j)\n"
	currentDayOfYearCurrentBody         = "# %j is a one-based, zero-padded calendar-year ordinal (001..366).\n${1:dayOfYear}=\\$(date +%j)\n"
	v6LocalShortBody                    = "${1:dateShort}=\\$(date -I) ${0:# format: yyyy/mm/dd}\n"
	currentLocalShortBody               = "${1:dateShort}=\\$(date '+%Y/%m/%d') ${0:# format: yyyy/mm/dd}\n"
	v6MonthNameCurrentBody              = "${1:monthName}=\\$(date +%${2|B,b|})\n"
	currentMonthNameCurrentBody         = "# Month names follow the current locale.\n${1:monthName}=\\$(date +%${2|B,b|})\n"
	v6MonthNumberCurrentBody            = "${1:monthNumber}=\\$(date +%m)\n"
	currentMonthNumberCurrentBody       = "# %m renders the one-based month number as 01..12.\n${1:monthNumber}=\\$(date +%m)\n"
	v6UTCLongBody                       = "${1:dateUTC}=\\$(date -u)\n"
	currentUTCLongBody                  = "# -u selects UTC; date's default display format follows LC_TIME.\n${1:dateUTC}=\\$(date -u)\n"
	v6YearCurrentBody                   = "${1:year}=\\$(date +%${2|Y,y|})\n"
	currentYearCurrentBody              = "# %Y and %y use the local calendar year, not the ISO week-based year.\n${1:year}=\\$(date +%${2|Y,y|})\n"
)

const v6FindFilesOrDirectoriesBody = "readarray -t resultArray < <(find \"/path/to/search\" -maxdepth 1 -type f -name \"criteria\")\n"
const currentFindFilesOrDirectoriesBody = "_shellman_find_root=\"/path/to/search\"\n_shellman_find_pattern=\"criteria\"\ncase $_shellman_find_root in\n  -*) _shellman_find_root=./$_shellman_find_root ;;\nesac\ncoproc _shellman_find {\n  find \"$_shellman_find_root\" -maxdepth 1 -type f -name \"$_shellman_find_pattern\" -print0\n}\n# shellcheck disable=SC2154 # Bash sets the named coprocess PID.\n_shellman_find_pid=$_shellman_find_PID\nreadarray -d '' -t resultArray <&\"${_shellman_find[0]}\"\nwait \"$_shellman_find_pid\"\n"

var v6IfDirectoryExistsBody = []any{
	`if [ -d "${1|/path/to/directory,${pathToDirectory}|}" ]; then`,
	"\t${2:echo \"directory \\\"${1}\\\" exists\"}",
	"fi\n",
}

var v6IfFileExecutableBody = []any{
	"if [ -x \"${1|/path/to/file,${filePath}|}\" ]; then",
	"\t${2:echo \"file \\\"${1}\\\" is executable\"}",
	"fi\n",
}

var v6IfFilesAreEqualBody = []any{
	`if cmp -s "${1|/path/to/file1,${filePath1}|}" "${2|/path/to/file2,${filePath2}|}"; then`,
	"\t${3:echo \"File \\\"${1}\\\" is equal to file \\\"${2}\\\"\"}",
	"fi\n",
}

const v6RemoveFilesOlderNewerBody = "find \"${1|/path/to/directory/,${pathToDirectory}|}\"* -type ${2|f,d|} -mtime ${3|-,+|}${4:days} -delete\n"

var currentIfFileExecutableBody = []any{
	"if (",
	"\t_shellman_executable_path=\"${1|/path/to/file,${filePath}|}\"",
	"\t[ -f \"$_shellman_executable_path\" ] && [ -x \"$_shellman_executable_path\" ]",
	"); then",
	"\t${2:echo \"file \\\"${1}\\\" is executable\"}",
	"fi\n",
}

var v6CtrlCBody = []any{
	"# CTRL+C event handler",
	"function on_ctrl_c() {",
	"\techo # Set cursor to the next line of '^C'",
	"\ttput cnorm # show cursor. You need this if animation is used.",
	"\t${1:# i.e. clean-up code here}",
	"\texit 1 # Don't remove. Use a number (1-255) for error code.",
	"}",
	"",
	"# Put this line at the beginning of your script (after functions used by event handlers).",
	"# Register CTRL+C event handler",
	"trap on_ctrl_c SIGINT",
}

var currentCtrlCBody = []any{
	"# CTRL+C event handler",
	"on_ctrl_c() {",
	"\techo # Set cursor to the next line of '^C'",
	"\ttput cnorm # show cursor. You need this if animation is used.",
	"\t${1:# i.e. clean-up code here}",
	"\texit 1 # Don't remove. Use a number (1-255) for error code.",
	"}",
	"",
	"# Put this line at the beginning of your script (after functions used by event handlers).",
	"# Register CTRL+C event handler",
	"trap on_ctrl_c INT",
}

var v6ExitBody = []any{
	"# Exit event handler",
	"function on_exit() {",
	"\ttput cnorm # Show cursor. You need this if animation is used.",
	"\t${1:# i.e. clean-up code here}",
	"\texit 0 # Exit gracefully.",
	"}",
	"",
	"# Put this line at the beginning of your script (after functions used by event handlers).",
	"# Register exit event handler.",
	"trap on_exit EXIT",
}

var currentExitBody = []any{
	"# Exit event handler",
	"on_exit() {",
	"\ttput cnorm # Show cursor. You need this if animation is used.",
	"\t${1:# i.e. clean-up code here}",
	"\texit 0 # Exit gracefully.",
	"}",
	"",
	"# Put this line at the beginning of your script (after functions used by event handlers).",
	"# Register exit event handler.",
	"trap on_exit EXIT",
}

var v6FileReadBody = []any{
	"while IFS= read -r line; do",
	"\techo \"\\$line\"",
	"done < ${1|/path/ro/file,\"${inputFile}\"|}\n",
}

var currentFileReadBody = []any{
	"file_path=${1|/path/ro/file,\"${inputFile}\"|}",
	"if [ -d \"\\$file_path\" ]; then",
	"\tprintf 'cannot read directory: %s\\n' \"\\$file_path\" >&2",
	"\tfalse",
	"else",
	"\twhile IFS= read -r line || [ -n \"\\$line\" ]; do",
	"\t\tprintf '%s\\n' \"\\$line\"",
	"\tdone < \"\\$file_path\"",
	"fi\n",
}

const v6FileSearchBody = "${1:result}=\\$(find \"${2|/path/to/search,${pathToSearch}|}\" -maxdepth ${3|1,2,3,4,5,6,7,8,9|} -type f -exec grep \"${4:criteria}\" {} +)\n"

const currentFileSearchBody = `${1:result}=\$(
  _shellman_search_path="${2|/path/to/search,${pathToSearch}|}"
  case \$_shellman_search_path in
    -*) _shellman_search_path=./\$_shellman_search_path ;;
  esac
  _shellman_status_file=\$(mktemp) || exit 2
  trap 'rm -f "\$_shellman_status_file"' 0
  find "\$_shellman_search_path" -maxdepth ${3|1,2,3,4,5,6,7,8,9|} -type f -exec sh -c '
    _shellman_pattern=\$1
    _shellman_status_file=\$2
    shift 2
    grep -e "\$_shellman_pattern" "\$@"
    _shellman_grep_status=\$?
    printf "%s\n" "\$_shellman_grep_status" >> "\$_shellman_status_file" || exit 2
  ' sh "${4:criteria}" "\$_shellman_status_file" {} +
  _shellman_find_status=\$?
  [ "\$_shellman_find_status" -eq 0 ] || exit 2
  _shellman_status_data=\$(cat "\$_shellman_status_file") || exit 2
  _shellman_match=1
  if [ -n "\$_shellman_status_data" ]; then
    while IFS= read -r _shellman_grep_status; do
      case \$_shellman_grep_status in
        0) _shellman_match=0 ;;
        1) ;;
        *) exit 2 ;;
      esac
    done <<_SHELLMAN_STATUS
\$_shellman_status_data
_SHELLMAN_STATUS
  fi
  exit "\$_shellman_match"
)
`

var v6FileWriteMultilineSudoBody = []any{
	`cat << EOF | sudo tee "${1|/path/to/file,${filePath}|}" >/dev/null`,
	"${2:first line\nsecond line\n...}",
	"EOF\n",
}

var currentFileWriteMultilineSudoBody = []any{
	`file_path="${1|/path/to/file,${filePath}|}"`,
	`case \$file_path in`,
	`  '') printf '%s\n' 'destination is empty' >&2; exit 2 ;;`,
	`  -*) file_path=./\$file_path ;;`,
	`esac`,
	`if [ ! -t 0 ]; then printf '%s\n' 'interactive terminal input is required' >&2; exit 2; fi`,
	`printf '%s\n' 'Enter text; press Ctrl-D at the start of a line to finish.' >&2`,
	"sudo tee \"\\$file_path\" >/dev/null || exit \\$?\n",
}

var v6FileWriteMultilineBody = []any{
	`cat >"${1|/path/to/file,${filePath}|}" <<EOF`,
	"${2:first line\nsecond line\n...}",
	"EOF\n",
}

var currentFileWriteMultilineBody = []any{
	`(`,
	`  file_path="${1|/path/to/file,${filePath}|}"`,
	`  case \$file_path in`,
	`    '') printf '%s\n' 'destination is empty' >&2; exit 2 ;;`,
	`    -*) file_path=./\$file_path ;;`,
	`  esac`,
	`  if [ ! -t 0 ]; then printf '%s\n' 'interactive terminal input is required' >&2; exit 2; fi`,
	`  printf '%s\n' 'Enter text; press Ctrl-D at the start of a line to finish.' >&2`,
	`  tee "\$file_path" >/dev/null`,
	")\n",
}

var v6FileWriteBody = []any{
	`echo "${1:sample header}" > "${2|/path/to/file,${pathToFile}|}"`,
	`echo # Empty line`,
	`for ${3:line} in "\${${4:lines}[@]}"; do`,
	"\techo \"\\${${3:line}\\}\" >> \"${2}\" # >> means append",
	"done\n",
}

var currentFileWriteBody = []any{
	`(`,
	`  file_path="${2|/path/to/file,${pathToFile}|}"`,
	`  case \$file_path in`,
	`    '') printf '%s\n' 'destination is empty' >&2; exit 2 ;;`,
	`    -*) file_path=./\$file_path ;;`,
	`  esac`,
	`  if [ ! -t 0 ]; then printf '%s\n' 'interactive terminal input is required' >&2; exit 2; fi`,
	`  printf '%s\n' 'Enter text; press Ctrl-D at the start of a line to finish.' >&2`,
	`  tee "\$file_path" >/dev/null`,
	")\n",
}

var v6FilesIterateBody = []any{
	"# Searching in the current directory",
	`for file in ${1|*,*.{jpg\,png}|}; do`,
	"\t${2:echo \"\\${file\\}\"}",
	"done\n",
}

var currentFilesIterateBody = []any{
	"# Searching in the current directory",
	"file_iterate_status=0",
	`for file in ${1|*,*.{jpg\,png}|}; do`,
	"\t[ -e \"\\$file\" ] || [ -L \"\\$file\" ] || continue",
	"\taction_status=0",
	"\t# shellcheck disable=SC2320",
	"\t${2:printf '%s\\n' \"\\${file\\}\"} || action_status=\\$?",
	"\tif [ \"\\$file_iterate_status\" -eq 0 ]; then file_iterate_status=\\$action_status; fi",
	"done",
	"(exit \"\\$file_iterate_status\")\n",
}

var v6DirectoriesIterateBody = []any{
	"# Make sure path ends with /",
	`for directory in ${1| /'path/to/directory'/,"${pathToDirectory}"|}*; do`,
	"\tif [[ -d \"\\${directory}\" && ! -L \"\\${directory}\" ]]; then",
	"\t\t${2:echo \"\\${directory\\}\"}",
	"\tfi",
	"done\n",
}

var currentDirectoriesIterateBody = []any{
	"# Make sure path ends with /",
	`for directory in ${1| /'path/to/directory'/,"${pathToDirectory}"|}*; do`,
	"\tif [ -d \"\\${directory}\" ] && [ ! -L \"\\${directory}\" ]; then",
	"\t\t${2:printf '%s\\n' \"\\${directory\\}\"}",
	"\tfi",
	"done\n",
}

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
	iteration68Ordered := append([]namedSnippet(nil), ordered...)
	iteration68JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-68 snippets", iteration68JSON, currentV7Iteration68SHA256)
	if ordered[67].name != "filesystem.remove-files-older-newer-than" {
		t.Fatalf("entry 68 = %q", ordered[67].name)
	}
	if ordered[67].snippet.Prefix != "remove old/new files/directories" || ordered[67].snippet.Description != "find and remove files(f)/directories(d) older(+)/newer(-) than x days" {
		t.Fatalf("entry-68 metadata = %#v", ordered[67].snippet)
	}
	ordered[67].snippet.Body = v6RemoveFilesOlderNewerBody
	changedByRollback68 := 0
	for index := range ordered {
		if !reflect.DeepEqual(iteration68Ordered[index].snippet, ordered[index].snippet) {
			changedByRollback68++
			if index != 67 {
				t.Fatalf("entry-68 rollback changed unexpected object %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if changedByRollback68 != 1 {
		t.Fatalf("entry-68 rollback changed %d objects, want 1", changedByRollback68)
	}
	rollback68JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-68 rollback to approved iteration 67", rollback68JSON, currentV7Iteration66SHA256)
	assertSHA256(t, "entry-68 rollback commands", renderDocumentation(ordered), historicalV6CommandsSHA256)
	iteration66Ordered := append([]namedSnippet(nil), ordered...)
	iteration66JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-66 snippets", iteration66JSON, currentV7Iteration66SHA256)
	if ordered[65].name != "filesystem.if-files-are-equal" {
		t.Fatalf("entry 66 = %q", ordered[65].name)
	}
	if ordered[65].snippet.Prefix != "if file =" || ordered[65].snippet.Description != "check if files are equal" {
		t.Fatalf("entry-66 metadata = %#v", ordered[65].snippet)
	}
	ordered[65].snippet.Body = append([]any(nil), v6IfFilesAreEqualBody...)
	changedByRollback66 := 0
	for index := range ordered {
		if !reflect.DeepEqual(iteration66Ordered[index].snippet, ordered[index].snippet) {
			changedByRollback66++
			if index != 65 {
				t.Fatalf("entry-66 rollback changed unexpected object %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if changedByRollback66 != 1 {
		t.Fatalf("entry-66 rollback changed %d objects, want 1", changedByRollback66)
	}
	rollback66JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-66 rollback to approved iteration 65", rollback66JSON, currentV7Iteration58SHA256)
	assertSHA256(t, "entry-66 rollback commands", renderDocumentation(ordered), historicalV6CommandsSHA256)
	iteration58JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-58 snippets", iteration58JSON, currentV7Iteration58SHA256)
	if ordered[57].name != "filesystem.if-file-executable" {
		t.Fatalf("entry 58 = %q", ordered[57].name)
	}
	if !reflect.DeepEqual(ordered[57].snippet.Body, currentIfFileExecutableBody) {
		t.Fatalf("entry-58 body = %#v", ordered[57].snippet.Body)
	}
	ordered[57].snippet.Body = append([]any(nil), v6IfFileExecutableBody...)
	rollback58JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-58 rollback to iteration 57", rollback58JSON, currentV7Iteration56SHA256)
	assertSHA256(t, "entry-58 rollback commands", renderDocumentation(ordered), historicalV6CommandsSHA256)
	if ordered[56].name != "filesystem.if-directory-exists" {
		t.Fatalf("entry 57 = %q", ordered[56].name)
	}
	if !reflect.DeepEqual(ordered[56].snippet.Body, v6IfDirectoryExistsBody) {
		t.Fatalf("entry-57 body = %#v", ordered[56].snippet.Body)
	}
	ordered[56].snippet.Body = append([]any(nil), v6IfDirectoryExistsBody...)
	rollback57JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-57 historical rollback to iteration 56", rollback57JSON, currentV7Iteration56SHA256)
	iteration56JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-56 snippets", iteration56JSON, currentV7Iteration56SHA256)
	if ordered[55].name != "filesystem.find-files-or-directories" {
		t.Fatalf("entry 56 = %q", ordered[55].name)
	}
	if ordered[55].snippet.Body != currentFindFilesOrDirectoriesBody {
		t.Fatal("entry-56 body differs")
	}
	ordered[55].snippet.Body = v6FindFilesOrDirectoriesBody
	rollback56JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-56 rollback to iteration 55", rollback56JSON, currentV7Iteration55SHA256)
	assertSHA256(t, "entry-56 rollback commands", renderDocumentation(ordered), historicalV6CommandsSHA256)
	iteration55Ordered := append([]namedSnippet(nil), ordered...)
	iteration55JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-55 snippets", iteration55JSON, currentV7Iteration55SHA256)
	if got := ordered[54].name; got != "filesystem.files-iterate" {
		t.Fatalf("traversal entry 55 = %q", got)
	}
	if !reflect.DeepEqual(ordered[54].snippet.Body, currentFilesIterateBody) {
		t.Fatalf("entry-55 body = %#v", ordered[54].snippet.Body)
	}
	ordered[54].snippet.Body = v6FilesIterateBody
	rollback55Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(iteration55Ordered[index].snippet, ordered[index].snippet) {
			rollback55Differences++
			if index != 54 {
				t.Fatalf("entry-55 rollback changed unexpected object %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback55Differences != 1 {
		t.Fatalf("entry-55 rollback changed %d objects", rollback55Differences)
	}
	rollback55JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-55 rollback to iteration 54", rollback55JSON, currentV7Iteration54SHA256)
	iteration54Ordered := append([]namedSnippet(nil), ordered...)
	iteration54JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-54 snippets", iteration54JSON, currentV7Iteration54SHA256)
	if got := ordered[53].name; got != "filesystem.file-write" {
		t.Fatalf("traversal entry 54 = %q", got)
	}
	if !reflect.DeepEqual(ordered[53].snippet.Body, currentFileWriteBody) {
		t.Fatalf("entry-54 body = %#v", ordered[53].snippet.Body)
	}
	ordered[53].snippet.Body = v6FileWriteBody
	rollback54Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(iteration54Ordered[index].snippet, ordered[index].snippet) {
			rollback54Differences++
			if index != 53 {
				t.Fatalf("entry-54 rollback changed unexpected object %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback54Differences != 1 {
		t.Fatalf("entry-54 rollback changed %d objects", rollback54Differences)
	}
	rollback54JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-54 rollback to iteration 53", rollback54JSON, currentV7Iteration53SHA256)
	iteration53Ordered := append([]namedSnippet(nil), ordered...)
	iteration53JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-53 snippets", iteration53JSON, currentV7Iteration53SHA256)
	if got := ordered[52].name; got != "filesystem.file-write-multiline" {
		t.Fatalf("traversal entry 53 = %q", got)
	}
	if !reflect.DeepEqual(ordered[52].snippet.Body, currentFileWriteMultilineBody) {
		t.Fatalf("entry-53 body = %#v", ordered[52].snippet.Body)
	}
	ordered[52].snippet.Body = v6FileWriteMultilineBody
	rollback53Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(iteration53Ordered[index].snippet, ordered[index].snippet) {
			rollback53Differences++
			if index != 52 {
				t.Fatalf("entry-53 rollback changed unexpected object %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback53Differences != 1 {
		t.Fatalf("entry-53 rollback changed %d objects", rollback53Differences)
	}
	rollback53JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-53 rollback to iteration 52", rollback53JSON, currentV7Iteration52SHA256)
	currentOrdered := append([]namedSnippet(nil), ordered...)
	iteration52JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-52 snippets", iteration52JSON, currentV7Iteration52SHA256)
	if got := ordered[51].name; got != "filesystem.file-write-multiline-sudo" {
		t.Fatalf("traversal entry 52 = %q", got)
	}
	if got := ordered[52].name; got != "filesystem.file-write-multiline" {
		t.Fatalf("traversal entry 53 = %q", got)
	}
	if !reflect.DeepEqual(ordered[51].snippet.Body, currentFileWriteMultilineSudoBody) {
		t.Fatalf("entry-52 body = %#v", ordered[51].snippet.Body)
	}
	ordered[51].snippet.Body = v6FileWriteMultilineSudoBody
	rollback52Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(currentOrdered[index].snippet, ordered[index].snippet) {
			rollback52Differences++
			if index != 51 {
				t.Fatalf("entry-52 rollback changed unexpected object %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback52Differences != 1 {
		t.Fatalf("entry-52 rollback changed %d objects", rollback52Differences)
	}
	rollback52JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-52 rollback to iteration 51", rollback52JSON, currentV7Iteration51SHA256)
	currentOrdered = append([]namedSnippet(nil), ordered...)
	iteration51JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-51 snippets", iteration51JSON, currentV7Iteration51SHA256)
	if got := ordered[50].name; got != "filesystem.file-search" {
		t.Fatalf("traversal entry 51 = %q", got)
	}
	if got := ordered[51].name; got != "filesystem.file-write-multiline-sudo" {
		t.Fatalf("traversal entry 52 = %q", got)
	}
	if got := ordered[50].snippet.Body; got != currentFileSearchBody {
		t.Fatalf("entry-51 body = %#v", got)
	}
	ordered[50].snippet.Body = v6FileSearchBody
	rollback51Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(currentOrdered[index].snippet, ordered[index].snippet) {
			rollback51Differences++
			if index != 50 {
				t.Fatalf("entry-51 rollback changed unexpected entry %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback51Differences != 1 {
		t.Fatalf("entry-51 rollback changed %d objects", rollback51Differences)
	}
	rollback51JSON, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "entry-51 rollback to iteration 50", rollback51JSON, currentV7Iteration50SHA256)
	currentOrdered = append([]namedSnippet(nil), ordered...)
	if got := ordered[44].name; got != iteration45MigrationName {
		t.Fatalf("traversal entry 45 = %q", got)
	}
	if got := ordered[45].name; got != iteration46MigrationName {
		t.Fatalf("traversal entry 46 = %q", got)
	}
	if got := ordered[46].name; got != iteration47MigrationName {
		t.Fatalf("traversal entry 47 = %q", got)
	}
	if got := ordered[47].name; got != iteration48MigrationName {
		t.Fatalf("traversal entry 48 = %q", got)
	}
	if got := ordered[48].name; got != iteration49MigrationName {
		t.Fatalf("traversal entry 49 = %q", got)
	}
	if got := ordered[49].name; got != iteration50MigrationName {
		t.Fatalf("traversal entry 50 = %q", got)
	}
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
	if got := ordered[34].name; got != iteration35MigrationName {
		t.Fatalf("thirty-fifth snippet in actual nsroot traversal = %q, want %q", got, iteration35MigrationName)
	}
	if got := ordered[35].name; got != iteration36MigrationName {
		t.Fatalf("thirty-sixth snippet in actual nsroot traversal = %q, want %q", got, iteration36MigrationName)
	}
	if got := ordered[36].name; got != iteration37MigrationName {
		t.Fatalf("thirty-seventh snippet in actual nsroot traversal = %q, want %q", got, iteration37MigrationName)
	}
	if got := ordered[37].name; got != iteration38MigrationName {
		t.Fatalf("thirty-eighth snippet in actual nsroot traversal = %q, want %q", got, iteration38MigrationName)
	}
	if got := ordered[38].name; got != iteration39MigrationName {
		t.Fatalf("thirty-ninth snippet in actual nsroot traversal = %q, want %q", got, iteration39MigrationName)
	}
	if got := ordered[39].name; got != iteration40MigrationName {
		t.Fatalf("fortieth snippet in actual nsroot traversal = %q, want %q", got, iteration40MigrationName)
	}
	if got := ordered[40].name; got != iteration41MigrationName {
		t.Fatalf("forty-first snippet in actual nsroot traversal = %q, want %q", got, iteration41MigrationName)
	}
	if got := ordered[41].name; got != iteration42MigrationName {
		t.Fatalf("forty-second snippet in actual nsroot traversal = %q, want %q", got, iteration42MigrationName)
	}
	if got := ordered[42].name; got != iteration43MigrationName {
		t.Fatalf("forty-third snippet in actual nsroot traversal = %q, want %q", got, iteration43MigrationName)
	}
	if got := ordered[43].name; got != iteration44MigrationName {
		t.Fatalf("forty-fourth snippet in actual nsroot traversal = %q, want %q", got, iteration44MigrationName)
	}
	if got := ordered[44].name; got != iteration45MigrationName {
		t.Fatalf("forty-fifth snippet in actual nsroot traversal = %q, want %q", got, iteration45MigrationName)
	}
	if got := ordered[46].snippet.Body; got != currentDirectoryCreateBody {
		t.Fatalf("entry-47 body = %#v", got)
	}
	if got := ordered[47].snippet.Body; got != currentDirectoryDeleteNestedBody {
		t.Fatalf("entry-48 body = %#v", got)
	}
	if got := ordered[48].snippet.Body; got != currentFileDeleteBody {
		t.Fatalf("entry-49 body = %#v", got)
	}
	if got := ordered[49].snippet.Body; !reflect.DeepEqual(got, currentFileReadBody) {
		t.Fatalf("entry-50 body = %#v", got)
	}
	currentIteration50, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-50 snippets", currentIteration50, currentV7Iteration50SHA256)
	ordered[49].snippet.Body = v6FileReadBody
	rollback50Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(currentOrdered[index].snippet, ordered[index].snippet) {
			rollback50Differences++
			if index != 49 || ordered[index].name != iteration50MigrationName {
				t.Fatalf("entry-50 rollback changed unexpected entry %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback50Differences != 1 {
		t.Fatalf("entry-50 rollback changed %d objects", rollback50Differences)
	}
	currentIteration49, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-49 snippets", currentIteration49, currentV7Iteration49SHA256)
	preRollback49 := append([]namedSnippet(nil), ordered...)
	ordered[48].snippet.Body = v6FileDeleteBody
	rollback49Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback49[index].snippet, ordered[index].snippet) {
			rollback49Differences++
			if index != 48 || ordered[index].name != iteration49MigrationName {
				t.Fatalf("entry-49 rollback changed unexpected entry %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback49Differences != 1 {
		t.Fatalf("entry-49 rollback changed %d objects", rollback49Differences)
	}
	currentIteration48, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-48 snippets", currentIteration48, currentV7Iteration48SHA256)
	preRollback48 := append([]namedSnippet(nil), ordered...)
	ordered[47].snippet.Body = v6DirectoryDeleteNestedBody
	rollback48Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback48[index].snippet, ordered[index].snippet) {
			rollback48Differences++
			if index != 47 || ordered[index].name != iteration48MigrationName {
				t.Fatalf("entry-48 rollback changed unexpected entry %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback48Differences != 1 {
		t.Fatalf("entry-48 rollback changed %d objects", rollback48Differences)
	}
	currentIteration47, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-47 snippets", currentIteration47, currentV7Iteration47SHA256)
	preRollback47 := append([]namedSnippet(nil), ordered...)
	ordered[46].snippet.Body = v6DirectoryCreateBody
	rollback47Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback47[index].snippet, ordered[index].snippet) {
			rollback47Differences++
			if index != 46 || ordered[index].name != iteration47MigrationName {
				t.Fatalf("entry-47 rollback changed unexpected entry %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback47Differences != 1 {
		t.Fatalf("entry-47 rollback changed %d objects", rollback47Differences)
	}
	currentIteration46, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "current v7 iteration-46 snippets", currentIteration46, currentV7Iteration46SHA256)
	assertSHA256(t, "current v7 iteration-46 commands", renderDocumentation(ordered), historicalV6CommandsSHA256)

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
	foundThirtyFifth := false
	foundThirtySixth := false
	foundThirtySeventh := false
	foundThirtyEighth := false
	foundThirtyNinth := false
	foundFortieth := false
	foundFortyFirst := false
	foundFortySecond := false
	foundFortyThird := false
	foundFortyFourth := false
	foundFortyFifth := false
	if got := ordered[45].snippet.Body; got != currentDirectoryCreateNestedBody {
		t.Fatalf("iteration-46 body = %#v", got)
	}
	preRollback46 := append([]namedSnippet(nil), ordered...)
	ordered[45].snippet.Body = v6DirectoryCreateNestedBody
	rollback46Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback46[index].snippet, ordered[index].snippet) {
			rollback46Differences++
			if index != 45 || ordered[index].name != iteration46MigrationName {
				t.Fatalf("iteration-46 rollback changed unexpected entry %d: %q", index+1, ordered[index].name)
			}
		}
	}
	if rollback46Differences != 1 {
		t.Fatalf("iteration-46 rollback changed %d objects", rollback46Differences)
	}
	reconstructedIteration45, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-45 snippets", reconstructedIteration45, currentV7Iteration45SHA256)
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
				t.Fatalf("iteration-34 date body = %#v, want %#v", got, currentDateNowShortBody)
			}
		case iteration35MigrationName:
			foundThirtyFifth = true
			if got, want := ordered[index].snippet.Prefix, "date now dayOfMonth"; got != want {
				t.Fatalf("iteration-35 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "current day of month (1..31)"; got != want {
				t.Fatalf("iteration-35 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDayOfMonthCurrentBody {
				t.Fatalf("iteration-35 body = %#v, want %#v", got, currentDayOfMonthCurrentBody)
			}
		case iteration36MigrationName:
			foundThirtySixth = true
			if got, want := ordered[index].snippet.Prefix, "date now dayOfWeek"; got != want {
				t.Fatalf("iteration-36 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "current day of week name (A: full, a: abbreviated)"; got != want {
				t.Fatalf("iteration-36 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDayOfWeekCurrentBody {
				t.Fatalf("iteration-36 body = %#v, want %#v", got, currentDayOfWeekCurrentBody)
			}
		case iteration37MigrationName:
			foundThirtySeventh = true
			if got, want := ordered[index].snippet.Prefix, "date now dayOfYear"; got != want {
				t.Fatalf("iteration-37 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "current day of year (1..366)"; got != want {
				t.Fatalf("iteration-37 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentDayOfYearCurrentBody {
				t.Fatalf("iteration-37 body = %#v, want %#v", got, currentDayOfYearCurrentBody)
			}
		case iteration38MigrationName:
			foundThirtyEighth = true
			if got, want := ordered[index].snippet.Prefix, "date now short"; got != want {
				t.Fatalf("iteration-38 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "yyyy/mm/dd"; got != want {
				t.Fatalf("iteration-38 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentLocalShortBody {
				t.Fatalf("iteration-38 body = %#v, want %#v", got, currentLocalShortBody)
			}
		case iteration39MigrationName:
			foundThirtyNinth = true
			if got, want := ordered[index].snippet.Prefix, "date now monthName"; got != want {
				t.Fatalf("iteration-39 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "current month name (B: full, b: abbreviated)"; got != want {
				t.Fatalf("iteration-39 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentMonthNameCurrentBody {
				t.Fatalf("iteration-39 body = %#v, want %#v", got, currentMonthNameCurrentBody)
			}
		case iteration40MigrationName:
			foundFortieth = true
			if got, want := ordered[index].snippet.Prefix, "date now monthNumber"; got != want {
				t.Fatalf("iteration-40 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "current month number (1..12)"; got != want {
				t.Fatalf("iteration-40 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentMonthNumberCurrentBody {
				t.Fatalf("iteration-40 body = %#v, want %#v", got, currentMonthNumberCurrentBody)
			}
		case iteration41MigrationName:
			foundFortyFirst = true
			if got, want := ordered[index].snippet.Prefix, "date now UTC"; got != want {
				t.Fatalf("iteration-41 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "coordinated Universal Time"; got != want {
				t.Fatalf("iteration-41 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentUTCLongBody {
				t.Fatalf("iteration-41 body = %#v, want %#v", got, currentUTCLongBody)
			}
		case iteration42MigrationName:
			foundFortySecond = true
			if got, want := ordered[index].snippet.Prefix, "date now year"; got != want {
				t.Fatalf("iteration-42 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "current Year (Y: full, y: last two digits)"; got != want {
				t.Fatalf("iteration-42 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; got != currentYearCurrentBody {
				t.Fatalf("iteration-42 body = %#v, want %#v", got, currentYearCurrentBody)
			}
		case iteration43MigrationName:
			foundFortyThird = true
			if got, want := ordered[index].snippet.Prefix, []any{"event CTRL+C", "event terminated"}; !reflect.DeepEqual(got, want) {
				t.Fatalf("iteration-43 prefixes = %#v, want %#v", got, want)
			}
			if got, want := ordered[index].snippet.Description, "register a function (handler) to run on script termination (CTRL+C)"; got != want {
				t.Fatalf("iteration-43 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentCtrlCBody) {
				t.Fatalf("iteration-43 body = %#v, want %#v", got, currentCtrlCBody)
			}
		case iteration44MigrationName:
			foundFortyFourth = true
			if got, want := ordered[index].snippet.Prefix, "event EXIT"; got != want {
				t.Fatalf("iteration-44 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "register a function (handler) to run on script exit"; got != want {
				t.Fatalf("iteration-44 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentExitBody) {
				t.Fatalf("iteration-44 body = %#v, want %#v", got, currentExitBody)
			}
		case iteration45MigrationName:
			foundFortyFifth = true
			if got, want := ordered[index].snippet.Prefix, "iterate directories"; got != want {
				t.Fatalf("iteration-45 prefix = %#v, want %q", got, want)
			}
			if got, want := ordered[index].snippet.Description, "write to a file"; got != want {
				t.Fatalf("iteration-45 description = %q, want %q", got, want)
			}
			if got := ordered[index].snippet.Body; !reflect.DeepEqual(got, currentDirectoriesIterateBody) {
				t.Fatalf("iteration-45 body = %#v, want %#v", got, currentDirectoriesIterateBody)
			}
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
	if !foundThirtyFifth {
		t.Fatalf("iteration-35 snippet %q is missing", iteration35MigrationName)
	}
	if !foundThirtySixth {
		t.Fatalf("iteration-36 snippet %q is missing", iteration36MigrationName)
	}
	if !foundThirtySeventh {
		t.Fatalf("iteration-37 snippet %q is missing", iteration37MigrationName)
	}
	if !foundThirtyEighth {
		t.Fatalf("iteration-38 snippet %q is missing", iteration38MigrationName)
	}
	if !foundThirtyNinth {
		t.Fatalf("iteration-39 snippet %q is missing", iteration39MigrationName)
	}
	if !foundFortieth {
		t.Fatalf("iteration-40 snippet %q is missing", iteration40MigrationName)
	}
	if !foundFortyFirst {
		t.Fatalf("iteration-41 snippet %q is missing", iteration41MigrationName)
	}
	if !foundFortySecond {
		t.Fatalf("iteration-42 snippet %q is missing", iteration42MigrationName)
	}
	if !foundFortyThird {
		t.Fatalf("iteration-43 snippet %q is missing", iteration43MigrationName)
	}
	if !foundFortyFourth {
		t.Fatalf("iteration-44 snippet %q is missing", iteration44MigrationName)
	}
	if !foundFortyFifth {
		t.Fatalf("iteration-45 snippet %q is missing", iteration45MigrationName)
	}
	preRollback45 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration45MigrationName {
			ordered[index].snippet.Body = v6DirectoriesIterateBody
		}
	}
	rollback45Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback45[index].snippet, ordered[index].snippet) {
			rollback45Differences++
			if ordered[index].name != iteration45MigrationName || index != 44 {
				t.Fatalf("iteration-45 rollback changed unexpected snippet %q at position %d", ordered[index].name, index+1)
			}
		}
	}
	if rollback45Differences != 1 {
		t.Fatalf("iteration-45 rollback changed %d snippets, want 1", rollback45Differences)
	}
	reconstructedIteration44, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-44 snippets", reconstructedIteration44, currentV7Iteration44SHA256)
	preRollback44 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration44MigrationName {
			ordered[index].snippet.Body = v6ExitBody
		}
	}
	rollback44Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback44[index].snippet, ordered[index].snippet) {
			rollback44Differences++
			if ordered[index].name != iteration44MigrationName || index != 43 {
				t.Fatalf("iteration-44 rollback changed unexpected snippet %q at position %d", ordered[index].name, index+1)
			}
		}
	}
	if rollback44Differences != 1 {
		t.Fatalf("iteration-44 rollback changed %d snippets, want 1", rollback44Differences)
	}
	reconstructedIteration43, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-43 snippets", reconstructedIteration43, currentV7Iteration43SHA256)
	preRollback43 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration43MigrationName {
			ordered[index].snippet.Body = v6CtrlCBody
		}
	}
	rollback43Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback43[index].snippet, ordered[index].snippet) {
			rollback43Differences++
			if ordered[index].name != iteration43MigrationName || index != 42 {
				t.Fatalf("iteration-43 rollback changed unexpected snippet %q at position %d", ordered[index].name, index+1)
			}
		}
	}
	if rollback43Differences != 1 {
		t.Fatalf("iteration-43 rollback changed %d snippets, want 1", rollback43Differences)
	}
	reconstructedIteration42, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-42 snippets", reconstructedIteration42, currentV7Iteration42SHA256)
	preRollback42 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration42MigrationName {
			ordered[index].snippet.Body = v6YearCurrentBody
		}
	}
	rollback42Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback42[index].snippet, ordered[index].snippet) {
			rollback42Differences++
			if ordered[index].name != iteration42MigrationName {
				t.Fatalf("iteration-42 rollback changed unexpected snippet %q", ordered[index].name)
			}
		}
	}
	if rollback42Differences != 1 {
		t.Fatalf("iteration-42 rollback changed %d snippets, want 1", rollback42Differences)
	}
	reconstructedIteration41, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-41 snippets", reconstructedIteration41, currentV7Iteration41SHA256)
	preRollback41 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration41MigrationName {
			ordered[index].snippet.Body = v6UTCLongBody
		}
	}
	rollback41Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback41[index].snippet, ordered[index].snippet) {
			rollback41Differences++
			if ordered[index].name != iteration41MigrationName {
				t.Fatalf("iteration-41 rollback changed unexpected snippet %q", ordered[index].name)
			}
		}
	}
	if rollback41Differences != 1 {
		t.Fatalf("iteration-41 rollback changed %d snippets, want 1", rollback41Differences)
	}
	reconstructedIteration40, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-40 snippets", reconstructedIteration40, currentV7Iteration40SHA256)
	preRollback40 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration40MigrationName {
			ordered[index].snippet.Body = v6MonthNumberCurrentBody
		}
	}
	rollback40Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback40[index].snippet, ordered[index].snippet) {
			rollback40Differences++
			if ordered[index].name != iteration40MigrationName {
				t.Fatalf("iteration-40 rollback changed unexpected snippet %q", ordered[index].name)
			}
		}
	}
	if rollback40Differences != 1 {
		t.Fatalf("iteration-40 rollback changed %d snippets, want 1", rollback40Differences)
	}
	reconstructedIteration39, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-39 snippets", reconstructedIteration39, currentV7Iteration39SHA256)
	preRollback39 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration39MigrationName {
			ordered[index].snippet.Body = v6MonthNameCurrentBody
		}
	}
	rollback39Differences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback39[index].snippet, ordered[index].snippet) {
			rollback39Differences++
			if ordered[index].name != iteration39MigrationName {
				t.Fatalf("iteration-39 rollback changed unexpected snippet %q", ordered[index].name)
			}
		}
	}
	if rollback39Differences != 1 {
		t.Fatalf("iteration-39 rollback changed %d snippets, want 1", rollback39Differences)
	}
	reconstructedIteration38, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-38 snippets", reconstructedIteration38, currentV7Iteration38SHA256)
	preRollback38 := append([]namedSnippet(nil), ordered...)
	for index := range ordered {
		if ordered[index].name == iteration38MigrationName {
			ordered[index].snippet.Body = v6LocalShortBody
		}
	}
	rollbackDifferences := 0
	for index := range ordered {
		if !reflect.DeepEqual(preRollback38[index].snippet, ordered[index].snippet) {
			rollbackDifferences++
			if ordered[index].name != iteration38MigrationName {
				t.Fatalf("iteration-38 rollback changed unexpected snippet %q", ordered[index].name)
			}
		}
	}
	if rollbackDifferences != 1 {
		t.Fatalf("iteration-38 rollback changed %d snippets, want 1", rollbackDifferences)
	}
	reconstructedIteration37, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-37 snippets", reconstructedIteration37, currentV7Iteration37SHA256)
	for index := range ordered {
		if ordered[index].name == iteration37MigrationName {
			ordered[index].snippet.Body = v6DayOfYearCurrentBody
		}
	}
	reconstructedIteration36, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-36 snippets", reconstructedIteration36, currentV7Iteration36SHA256)
	for index := range ordered {
		if ordered[index].name == iteration36MigrationName {
			ordered[index].snippet.Body = v6DayOfWeekCurrentBody
		}
	}
	reconstructedIteration35, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-35 snippets", reconstructedIteration35, currentV7Iteration35SHA256)
	for index := range ordered {
		if ordered[index].name == iteration35MigrationName {
			ordered[index].snippet.Body = v6DayOfMonthCurrentBody
		}
	}
	reconstructedIteration34, err := renderSnippetJSON(ordered)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "reconstructed approved v7 iteration-34 snippets", reconstructedIteration34, currentV7Iteration34SHA256)
	for index := range ordered {
		if ordered[index].name == iteration34MigrationName {
			ordered[index].snippet.Body = v6DateNowShortBody
		}
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
		if !reflect.DeepEqual(iteration68Ordered[index].snippet, ordered[index].snippet) {
			changed++
		}
	}
	if changed != 59 {
		t.Fatalf("snippets differing from historical v6 = %d, want 59 (219 unchanged)", changed)
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

func expandedGeneratedDateNowShort(t *testing.T, destination string) string {
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
		t.Fatalf("generated date body = %#v, want %#v", snippets[iteration34MigrationName].Body, currentDateNowShortBody)
	}
	if strings.Count(body, `${1:dateShort}`) != 1 || strings.Count(body, `${0:# format: yyyy/mm/dd}`) != 1 ||
		strings.Count(body, `\$`) != 1 || !strings.HasSuffix(body, "\n") {
		t.Fatalf("unexpected date placeholder or escaping structure: %q", body)
	}
	body = strings.ReplaceAll(body, `${1:dateShort}`, destination)
	body = strings.ReplaceAll(body, `${0:# format: yyyy/mm/dd}`, `# format: yyyy/mm/dd`)
	return strings.ReplaceAll(body, `\$`, `$`)
}

func expandedGeneratedMonthName(t *testing.T, destination, choice string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration39MigrationName].Body.(string)
	if !ok || body != currentMonthNameCurrentBody {
		t.Fatalf("generated month body = %#v", snippets[iteration39MigrationName].Body)
	}
	if strings.Count(body, `${1:monthName}`) != 1 || strings.Count(body, `${2|B,b|}`) != 1 ||
		strings.Count(body, `\$`) != 1 || !strings.HasSuffix(body, "\n") {
		t.Fatalf("unexpected month placeholder structure: %q", body)
	}
	if choice != "B" && choice != "b" {
		t.Fatalf("invalid month choice %q", choice)
	}
	body = strings.ReplaceAll(body, `${1:monthName}`, destination)
	body = strings.ReplaceAll(body, `${2|B,b|}`, choice)
	return strings.ReplaceAll(body, `\$`, `$`)
}

func TestMonthNameCurrentGeneratedExpansionAndBehavior(t *testing.T) {
	requireCommand(t, "sh")
	for _, choice := range []string{"B", "b"} {
		got := expandedGeneratedMonthName(t, "month_result", choice)
		want := "# Month names follow the current locale.\nmonth_result=$(date +%" + choice + ")\n"
		if got != want {
			t.Fatalf("choice %s: expansion = %q, want %q", choice, got, want)
		}
		path := filepath.Join(t.TempDir(), "month.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+got+"printf '%s\\n' \"$month_result\" >/dev/null\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "shellcheck", "-s", "sh", path)
		if err != nil {
			t.Fatalf("ShellCheck choice %s: %v: %s", choice, err, output)
		}
	}

	// Only the test stand-in uses GNU date -d; the snippet invokes date +%B or date +%b.
	full := []string{"January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"}
	abbreviated := []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}
	for month := 1; month <= 12; month++ {
		for _, choice := range []string{"B", "b"} {
			instant := fmt.Sprintf("2024-%02d-15 12:00:00 UTC", month)
			want := full[month-1]
			if choice == "b" {
				want = abbreviated[month-1]
			}
			t.Run(fmt.Sprintf("month-%02d-%s", month, choice), func(t *testing.T) {
				root := t.TempDir()
				bin := filepath.Join(root, "bin")
				mustMkdirAll(t, bin)
				mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = '+%"+choice+"' ] || exit 42\nexec /bin/date -d \"$SHELLMAN_MONTH_INSTANT\" \"$1\"\n")
				path := filepath.Join(root, "month.sh")
				mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedMonthName(t, "month_result", choice)+"[ \"$month_result\" = '"+want+"' ]\n"))
				output, err := runCommand(".", append(os.Environ(), "PATH="+bin, "LC_ALL=C", "TZ=UTC", "SHELLMAN_MONTH_INSTANT="+instant), "sh", path)
				if err != nil {
					t.Fatalf("want %q: %v: %s", want, err, output)
				}
			})
		}
	}
}

func TestMonthNameCurrentAssignmentStatusAndText(t *testing.T) {
	requireCommand(t, "sh")
	runCase := func(t *testing.T, choice, stub, setup, destination, trailer string, includeDate bool) (string, int) {
		t.Helper()
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		if includeDate {
			mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n"+stub)
		}
		path := filepath.Join(root, "month.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+expandedGeneratedMonthName(t, destination, choice)+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		return string(output), exitCode(err)
	}
	for _, choice := range []string{"B", "b"} {
		t.Run("arbitrary-text-"+choice, func(t *testing.T) {
			stub := "[ \"$#\" -eq 1 ] && [ \"$1\" = '+%" + choice + "' ] || exit 42\nprintf '\\nMois. long\\ninterne\\n\\n'\n"
			output, status := runCase(t, choice, stub, "", "month_result", `[ "$snippet_status" -eq 0 ] && [ "$month_result" = "$(printf '\nMois. long\ninterne')" ]`, true)
			if status != 0 || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
		for _, failure := range []int{1, 42, 127} {
			t.Run(fmt.Sprintf("failure-%s-%d", choice, failure), func(t *testing.T) {
				stub := fmt.Sprintf("printf 'partial\\n'; printf 'month error\\n' >&2; exit %d\n", failure)
				output, status := runCase(t, choice, stub, "", "month_result", `[ "$month_result" = partial ] && exit "$snippet_status"`, true)
				if status != failure || !strings.Contains(output, "month error") {
					t.Fatalf("status=%d output=%q", status, output)
				}
			})
		}
	}
	t.Run("command-not-found", func(t *testing.T) {
		output, status := runCase(t, "B", "", "", "month_result", `exit "$snippet_status"`, false)
		if status != 127 || !strings.Contains(output, "date") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	for _, destination := range []string{"month_result", "1bad"} {
		t.Run("destination-"+destination, func(t *testing.T) {
			setup := ""
			if destination == "month_result" {
				setup = "readonly month_result"
			}
			output, status := runCase(t, "B", "printf 'March\\n'\n", setup, destination, `exit "$snippet_status"`, true)
			if status == 0 || !strings.Contains(output, destination) {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}
	t.Run("later-command-masks-status", func(t *testing.T) {
		output, status := runCase(t, "B", "exit 42\n", "", "month_result", `:; [ "$snippet_status" -eq 42 ]`, true)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("pipeline-masks-left-status", func(t *testing.T) {
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\nexit 42\n")
		path := filepath.Join(root, "pipeline.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n{\n"+expandedGeneratedMonthName(t, "month_result", "B")+"} | /bin/cat\n"))
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		if status := exitCode(err); status != 0 || string(output) != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func expandedGeneratedMonthNumber(t *testing.T, destination string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration40MigrationName].Body.(string)
	if !ok || body != currentMonthNumberCurrentBody {
		t.Fatalf("generated month-number body = %#v", snippets[iteration40MigrationName].Body)
	}
	if strings.Count(body, `${1:monthNumber}`) != 1 || strings.Count(body, `\$`) != 1 ||
		strings.Count(body, `date +%m`) != 1 || !strings.HasSuffix(body, "\n") {
		t.Fatalf("unexpected month-number placeholder structure: %q", body)
	}
	body = strings.ReplaceAll(body, `${1:monthNumber}`, destination)
	return strings.ReplaceAll(body, `\$`, `$`)
}

func TestMonthNumberCurrentGeneratedExpansionAndMonths(t *testing.T) {
	requireCommand(t, "sh")
	expanded := expandedGeneratedMonthNumber(t, "month_result")
	if want := "# %m renders the one-based month number as 01..12.\nmonth_result=$(date +%m)\n"; expanded != want {
		t.Fatalf("expanded body = %q, want %q", expanded, want)
	}
	path := filepath.Join(t.TempDir(), "month-number.sh")
	mustWriteFile(t, path, []byte("#!/bin/sh\n"+expanded+"printf '%s\\n' \"$month_result\" >/dev/null\n"))
	run(t, ".", "sh", "-n", path)
	output, err := runCommand(".", nil, "shellcheck", "-s", "sh", path)
	if err != nil {
		t.Fatalf("ShellCheck: %v: %s", err, output)
	}

	// GNU date -d controls test time only; the snippet itself passes just +%m.
	for month := 1; month <= 12; month++ {
		t.Run(fmt.Sprintf("month-%02d", month), func(t *testing.T) {
			root := t.TempDir()
			bin := filepath.Join(root, "bin")
			mustMkdirAll(t, bin)
			mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = '+%m' ] || exit 42\nexec /bin/date -d \"$SHELLMAN_MONTH_INSTANT\" \"$1\"\n")
			path := filepath.Join(root, "month.sh")
			want := fmt.Sprintf("%02d", month)
			mustWriteFile(t, path, []byte("#!/bin/sh\n"+expanded+"[ \"$month_result\" = '"+want+"' ]\n"))
			instant := fmt.Sprintf("2024-%02d-15 12:00:00 UTC", month)
			output, err := runCommand(".", append(os.Environ(), "PATH="+bin, "LC_ALL=C", "TZ=UTC", "SHELLMAN_MONTH_INSTANT="+instant), "sh", path)
			if err != nil {
				t.Fatalf("month %d, want %q: %v: %s", month, want, err, output)
			}
		})
	}
}

func TestMonthNumberCurrentAssignmentStatusAndText(t *testing.T) {
	requireCommand(t, "sh")
	runCase := func(t *testing.T, stub, setup, destination, trailer string, includeDate bool) (string, int) {
		t.Helper()
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		if includeDate {
			mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n"+stub)
		}
		path := filepath.Join(root, "month.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+expandedGeneratedMonthNumber(t, destination)+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		return string(output), exitCode(err)
	}
	t.Run("arbitrary-text-and-newlines", func(t *testing.T) {
		stub := "[ \"$#\" -eq 1 ] && [ \"$1\" = '+%m' ] || exit 42\nprintf '\\n0 7\\ninner\\n\\n'\n"
		output, status := runCase(t, stub, "", "month_result", `[ "$snippet_status" -eq 0 ] && [ "$month_result" = "$(printf '\n0 7\ninner')" ]`, true)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	for _, failure := range []int{1, 42, 127} {
		t.Run(fmt.Sprintf("failure-%d", failure), func(t *testing.T) {
			stub := fmt.Sprintf("printf 'partial\\n'; printf 'date error\\n' >&2; exit %d\n", failure)
			output, status := runCase(t, stub, "", "month_result", `[ "$month_result" = partial ] && exit "$snippet_status"`, true)
			if status != failure || !strings.Contains(output, "date error") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}
	t.Run("command-not-found", func(t *testing.T) {
		output, status := runCase(t, "", "", "month_result", `exit "$snippet_status"`, false)
		if status != 127 || !strings.Contains(output, "date") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	for _, destination := range []string{"month_result", "1bad"} {
		t.Run("destination-"+destination, func(t *testing.T) {
			setup := ""
			if destination == "month_result" {
				setup = "readonly month_result"
			}
			output, status := runCase(t, "printf '04\\n'\n", setup, destination, `exit "$snippet_status"`, true)
			if status == 0 || !strings.Contains(output, destination) {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}
	t.Run("later-command-masks-status", func(t *testing.T) {
		output, status := runCase(t, "exit 42\n", "", "month_result", `:; [ "$snippet_status" -eq 42 ]`, true)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("pipeline-masks-left-status", func(t *testing.T) {
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\nexit 42\n")
		path := filepath.Join(root, "pipeline.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n{\n"+expandedGeneratedMonthNumber(t, "month_result")+"} | /bin/cat\n"))
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		if status := exitCode(err); status != 0 || string(output) != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func expandedGeneratedUTCLong(t *testing.T, destination string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration41MigrationName].Body.(string)
	if !ok || body != currentUTCLongBody {
		t.Fatalf("generated UTC body = %#v", snippets[iteration41MigrationName].Body)
	}
	if strings.Count(body, `${1:dateUTC}`) != 1 || strings.Count(body, `\$`) != 1 ||
		strings.Count(body, `date -u`) != 1 || !strings.HasSuffix(body, "\n") {
		t.Fatalf("unexpected UTC placeholder structure: %q", body)
	}
	body = strings.ReplaceAll(body, `${1:dateUTC}`, destination)
	return strings.ReplaceAll(body, `\$`, `$`)
}

func TestUTCLongGeneratedExpansionAndBoundaries(t *testing.T) {
	requireCommand(t, "sh")
	expanded := expandedGeneratedUTCLong(t, "utc_result")
	if want := "# -u selects UTC; date's default display format follows LC_TIME.\nutc_result=$(date -u)\n"; expanded != want {
		t.Fatalf("expanded body = %q, want %q", expanded, want)
	}
	path := filepath.Join(t.TempDir(), "utc.sh")
	mustWriteFile(t, path, []byte("#!/bin/sh\n"+expanded+"printf '%s\\n' \"$utc_result\" >/dev/null\n"))
	run(t, ".", "sh", "-n", path)
	output, err := runCommand(".", nil, "shellcheck", "-s", "sh", path)
	if err != nil {
		t.Fatalf("ShellCheck: %v: %s", err, output)
	}

	// GNU date -d controls the test instant only; the generated snippet passes exactly -u.
	cases := []struct{ name, instant, want string }{
		{"year-end", "2024-12-31 23:59:59 UTC", "Tue Dec 31 23:59:59 UTC 2024"},
		{"new-year-midnight", "2025-01-01 00:00:00 UTC", "Wed Jan  1 00:00:00 UTC 2025"},
		{"month-end", "2024-01-31 23:59:59 UTC", "Wed Jan 31 23:59:59 UTC 2024"},
		{"new-month-midnight", "2024-02-01 00:00:00 UTC", "Thu Feb  1 00:00:00 UTC 2024"},
		{"leap-day-noon", "2024-02-29 12:00:00 UTC", "Thu Feb 29 12:00:00 UTC 2024"},
		{"after-leap-day", "2024-03-01 00:00:00 UTC", "Fri Mar  1 00:00:00 UTC 2024"},
		{"representative-daytime", "2024-06-15 12:34:56 UTC", "Sat Jun 15 12:34:56 UTC 2024"},
	}
	for _, test := range cases {
		for _, timezone := range []string{"UTC", "America/Los_Angeles", "Asia/Kathmandu", "Pacific/Kiritimati"} {
			t.Run(test.name+"-"+timezone, func(t *testing.T) {
				root := t.TempDir()
				bin := filepath.Join(root, "bin")
				mustMkdirAll(t, bin)
				mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = '-u' ] || exit 42\nexec /bin/date -u -d \"$SHELLMAN_UTC_INSTANT\"\n")
				path := filepath.Join(root, "utc.sh")
				mustWriteFile(t, path, []byte("#!/bin/sh\n"+expanded+"[ \"$utc_result\" = '"+test.want+"' ]\n"))
				output, err := runCommand(".", append(os.Environ(), "PATH="+bin, "LC_ALL=C", "TZ="+timezone, "SHELLMAN_UTC_INSTANT="+test.instant), "sh", path)
				if err != nil {
					t.Fatalf("want %q: %v: %s", test.want, err, output)
				}
			})
		}
	}
}

func TestUTCLongAssignmentStatusAndText(t *testing.T) {
	requireCommand(t, "sh")
	runCase := func(t *testing.T, stub, setup, destination, trailer string, includeDate bool) (string, int) {
		t.Helper()
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		if includeDate {
			mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n"+stub)
		}
		path := filepath.Join(root, "utc.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+expandedGeneratedUTCLong(t, destination)+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		return string(output), exitCode(err)
	}
	t.Run("arbitrary-text-and-newlines", func(t *testing.T) {
		stub := "[ \"$#\" -eq 1 ] && [ \"$1\" = '-u' ] || exit 42\nprintf '\\nUTC label\\ninner\\n\\n'\n"
		output, status := runCase(t, stub, "", "utc_result", `[ "$snippet_status" -eq 0 ] && [ "$utc_result" = "$(printf '\nUTC label\ninner')" ]`, true)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	for _, failure := range []int{1, 42, 127} {
		t.Run(fmt.Sprintf("failure-%d", failure), func(t *testing.T) {
			stub := fmt.Sprintf("printf 'partial\\n'; printf 'date error\\n' >&2; exit %d\n", failure)
			output, status := runCase(t, stub, "", "utc_result", `[ "$utc_result" = partial ] && exit "$snippet_status"`, true)
			if status != failure || !strings.Contains(output, "date error") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}
	t.Run("command-not-found", func(t *testing.T) {
		output, status := runCase(t, "", "", "utc_result", `exit "$snippet_status"`, false)
		if status != 127 || !strings.Contains(output, "date") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	for _, destination := range []string{"utc_result", "1bad"} {
		t.Run("destination-"+destination, func(t *testing.T) {
			setup := ""
			if destination == "utc_result" {
				setup = "readonly utc_result"
			}
			output, status := runCase(t, "printf 'UTC\\n'\n", setup, destination, `exit "$snippet_status"`, true)
			if status == 0 || !strings.Contains(output, destination) {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}
	t.Run("later-command-masks-status", func(t *testing.T) {
		output, status := runCase(t, "exit 42\n", "", "utc_result", `:; [ "$snippet_status" -eq 42 ]`, true)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("pipeline-masks-left-status", func(t *testing.T) {
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\nexit 42\n")
		path := filepath.Join(root, "pipeline.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n{\n"+expandedGeneratedUTCLong(t, "utc_result")+"} | /bin/cat\n"))
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		if status := exitCode(err); status != 0 || string(output) != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func expandedGeneratedYearCurrent(t *testing.T, destination, choice string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration42MigrationName].Body.(string)
	if !ok || body != currentYearCurrentBody {
		t.Fatalf("generated year body = %#v", snippets[iteration42MigrationName].Body)
	}
	if strings.Count(body, `${1:year}`) != 1 || strings.Count(body, `${2|Y,y|}`) != 1 ||
		strings.Count(body, `\$`) != 1 || !strings.HasSuffix(body, "\n") {
		t.Fatalf("unexpected year placeholder structure: %q", body)
	}
	if choice != "Y" && choice != "y" {
		t.Fatalf("invalid year choice %q", choice)
	}
	body = strings.ReplaceAll(body, `${1:year}`, destination)
	body = strings.ReplaceAll(body, `${2|Y,y|}`, choice)
	return strings.ReplaceAll(body, `\$`, `$`)
}

func TestYearCurrentGeneratedExpansionAndCalendarBoundaries(t *testing.T) {
	requireCommand(t, "sh")
	for _, choice := range []string{"Y", "y"} {
		expanded := expandedGeneratedYearCurrent(t, "year_result", choice)
		want := "# %Y and %y use the local calendar year, not the ISO week-based year.\nyear_result=$(date +%" + choice + ")\n"
		if expanded != want {
			t.Fatalf("choice %s: expansion = %q, want %q", choice, expanded, want)
		}
		path := filepath.Join(t.TempDir(), "year.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+expanded+"printf '%s\\n' \"$year_result\" >/dev/null\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", nil, "shellcheck", "-s", "sh", path)
		if err != nil {
			t.Fatalf("ShellCheck choice %s: %v: %s", choice, err, output)
		}
	}

	// GNU date -d fixes test instants only; the snippet passes exactly +%Y or +%y.
	cases := []struct{ name, instant, timezone, full, short string }{
		{"ordinary", "2026-06-15 12:00:00 UTC", "UTC", "2026", "26"},
		{"leap-day", "2024-02-29 12:00:00 UTC", "UTC", "2024", "24"},
		{"december-31", "2024-12-31 23:59:59 UTC", "UTC", "2024", "24"},
		{"january-1", "2025-01-01 00:00:00 UTC", "UTC", "2025", "25"},
		{"iso-next-year", "2019-12-30 12:00:00 UTC", "UTC", "2019", "19"},
		{"iso-previous-year", "2021-01-01 12:00:00 UTC", "UTC", "2021", "21"},
		{"same-instant-west", "2024-12-31 23:30:00 UTC", "America/Los_Angeles", "2024", "24"},
		{"same-instant-east", "2024-12-31 23:30:00 UTC", "Asia/Kathmandu", "2025", "25"},
	}
	for _, test := range cases {
		for _, choice := range []string{"Y", "y"} {
			t.Run(test.name+"-"+choice, func(t *testing.T) {
				root := t.TempDir()
				bin := filepath.Join(root, "bin")
				mustMkdirAll(t, bin)
				mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = '+%"+choice+"' ] || exit 42\nexec /bin/date -d \"$SHELLMAN_YEAR_INSTANT\" \"$1\"\n")
				path := filepath.Join(root, "year.sh")
				want := test.full
				if choice == "y" {
					want = test.short
				}
				mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedYearCurrent(t, "year_result", choice)+"[ \"$year_result\" = '"+want+"' ]\n"))
				output, err := runCommand(".", append(os.Environ(), "PATH="+bin, "LC_ALL=C", "TZ="+test.timezone, "SHELLMAN_YEAR_INSTANT="+test.instant), "sh", path)
				if err != nil {
					t.Fatalf("want %q: %v: %s", want, err, output)
				}
			})
		}
	}
}

func TestYearCurrentAssignmentStatusAndText(t *testing.T) {
	requireCommand(t, "sh")
	runCase := func(t *testing.T, choice, stub, setup, destination, trailer string, includeDate bool) (string, int) {
		t.Helper()
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		if includeDate {
			mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\n"+stub)
		}
		path := filepath.Join(root, "year.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+expandedGeneratedYearCurrent(t, destination, choice)+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		return string(output), exitCode(err)
	}
	for _, choice := range []string{"Y", "y"} {
		t.Run("arbitrary-text-"+choice, func(t *testing.T) {
			stub := "[ \"$#\" -eq 1 ] && [ \"$1\" = '+%" + choice + "' ] || exit 42\nprintf '\\n20 26\\ninside\\n\\n'\n"
			output, status := runCase(t, choice, stub, "", "year_result", `[ "$snippet_status" -eq 0 ] && [ "$year_result" = "$(printf '\n20 26\ninside')" ]`, true)
			if status != 0 || output != "" {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
		for _, failure := range []int{1, 42, 127} {
			t.Run(fmt.Sprintf("failure-%s-%d", choice, failure), func(t *testing.T) {
				stub := fmt.Sprintf("printf 'partial\\n'; printf 'date error\\n' >&2; exit %d\n", failure)
				output, status := runCase(t, choice, stub, "", "year_result", `[ "$year_result" = partial ] && exit "$snippet_status"`, true)
				if status != failure || !strings.Contains(output, "date error") {
					t.Fatalf("status=%d output=%q", status, output)
				}
			})
		}
	}
	t.Run("command-not-found", func(t *testing.T) {
		output, status := runCase(t, "Y", "", "", "year_result", `exit "$snippet_status"`, false)
		if status != 127 || !strings.Contains(output, "date") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	for _, destination := range []string{"year_result", "1bad"} {
		t.Run("destination-"+destination, func(t *testing.T) {
			setup := ""
			if destination == "year_result" {
				setup = "readonly year_result"
			}
			output, status := runCase(t, "Y", "printf '2026\\n'\n", setup, destination, `exit "$snippet_status"`, true)
			if status == 0 || !strings.Contains(output, destination) {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}
	t.Run("later-command-masks-status", func(t *testing.T) {
		output, status := runCase(t, "Y", "exit 42\n", "", "year_result", `:; [ "$snippet_status" -eq 42 ]`, true)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("pipeline-masks-left-status", func(t *testing.T) {
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		mustExecutable(t, filepath.Join(bin, "date"), "#!/bin/sh\nexit 42\n")
		path := filepath.Join(root, "pipeline.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n{\n"+expandedGeneratedYearCurrent(t, "year_result", "Y")+"} | /bin/cat\n"))
		output, err := runCommand(".", append(os.Environ(), "PATH="+bin), "sh", path)
		if status := exitCode(err); status != 0 || string(output) != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func expandedGeneratedCtrlC43(t *testing.T, action string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration43MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentCtrlCBody) {
		t.Fatalf("generated entry-43 body = %#v", snippets[iteration43MigrationName].Body)
	}
	lines := make([]string, len(body))
	for i, value := range body {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d is not a string", i)
		}
		lines[i] = line
	}
	joined := strings.Join(lines, "\n")
	placeholder := `${1:# i.e. clean-up code here}`
	if strings.Count(joined, placeholder) != 1 || strings.Contains(joined, "${0") ||
		!strings.HasSuffix(joined, "trap on_ctrl_c INT") {
		t.Fatalf("unexpected generated placeholder or trap: %q", joined)
	}
	return strings.Replace(joined, placeholder, action, 1)
}

func TestCtrlC43GeneratedSyntaxAndIsolatedSignal(t *testing.T) {
	for _, shell := range []string{"sh", "bash", "zsh"} {
		requireCommand(t, shell)
		path := filepath.Join(t.TempDir(), "syntax.sh")
		mustWriteFile(t, path, []byte(expandedGeneratedCtrlC43(t, "# i.e. clean-up code here")+"\n"))
		run(t, ".", shell, "-n", path)
	}
	defaultPath := filepath.Join(t.TempDir(), "default.sh")
	defaultBody := expandedGeneratedCtrlC43(t, "# i.e. clean-up code here")
	if !strings.Contains(defaultBody, "on_ctrl_c() {\n") || !strings.Contains(defaultBody, "\n\texit 1 # Don't remove") {
		t.Fatalf("unexpected exact default expansion: %q", defaultBody)
	}
	mustWriteFile(t, defaultPath, []byte("#!/bin/sh\n"+defaultBody+"\n"))
	output, err := runCommand(".", nil, "shellcheck", "-s", "sh", defaultPath)
	if err != nil {
		t.Fatalf("ShellCheck on generated default: %v: %s", err, output)
	}

	type signalCase struct {
		name, shell, action, setup, after, wantTrace, wantOutput string
		wantStatus                                               int
		fakeTput                                                 bool
	}
	cases := []signalCase{
		{name: "default shell source", shell: "sh", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "TPUT\nACTION\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "late expansion and prior trap replacement", shell: "sh", setup: `value=before; trap 'printf "OLD\n" >> "$TRACE"' INT`, after: `value=after`, action: `printf '%s\n' "$(printf 'SUB\n' >> "$TRACE"; printf '%s' "$value")" >> "$TRACE"; printf 'SECOND\n' >> "$TRACE"`, wantTrace: "TPUT\nSUB\nafter\nSECOND\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "quoted spaces", shell: "sh", action: `printf '%s\n' 'two words' >> "$TRACE"`, wantTrace: "TPUT\ntwo words\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "parameter arithmetic and status at action", shell: "sh", setup: `value=before`, after: `value=after`, action: `printf 'STATUS:%s ARITH:%s VAR:%s\n' "$?" "$((2+3))" "${value:-unset}" >> "$TRACE"`, wantTrace: "TPUT\nSTATUS:0 ARITH:5 VAR:after\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "handler failure masked by explicit exit", shell: "sh", action: `false; printf 'AFTER_FAIL\n' >> "$TRACE"`, wantTrace: "TPUT\nAFTER_FAIL\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "explicit exit overrides", shell: "sh", action: `exit 42`, wantTrace: "TPUT\n", wantOutput: "\n", wantStatus: 42, fakeTput: true},
		{name: "explicit successful exit overrides", shell: "sh", action: `exit 0`, wantTrace: "TPUT\n", wantOutput: "\n", wantStatus: 0, fakeTput: true},
		{name: "empty source placeholder", shell: "sh", action: ``, wantTrace: "TPUT\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "unavailable action command", shell: "sh", action: `missing_shellman_command; printf 'CONTINUED\n' >> "$TRACE"`, wantTrace: "TPUT\nCONTINUED\n", wantOutput: "", wantStatus: 1, fakeTput: true},
		{name: "unavailable tput", shell: "sh", action: `printf 'TPUT_STATUS:%s\n' "$?" >> "$TRACE"`, wantTrace: "TPUT_STATUS:127\n", wantOutput: "", wantStatus: 1},
		{name: "bash runtime", shell: "bash", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "TPUT\nACTION\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
		{name: "zsh runtime", shell: "zsh", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "TPUT\nACTION\n", wantOutput: "\n", wantStatus: 1, fakeTput: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			bin := filepath.Join(root, "bin")
			mustMkdirAll(t, bin)
			if tc.fakeTput {
				mustExecutable(t, filepath.Join(bin, "tput"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = cnorm ] || exit 42\nprintf 'TPUT\\n' >> \"$TRACE\"\n")
			}
			trace, ready := filepath.Join(root, "trace"), filepath.Join(root, "ready")
			script := filepath.Join(root, "child.sh")
			body := expandedGeneratedCtrlC43(t, tc.action)
			mustWriteFile(t, script, []byte("#!/bin/sh\n"+tc.setup+"\n"+body+"\n"+tc.after+"\nprintf 'ready\\n' > \"$READY\"\nwhile :; do :; done\n"))
			run(t, ".", tc.shell, "-n", script)
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, tc.shell, script)
			cmd.Env = append(os.Environ(), "PATH="+bin, "TRACE="+trace, "READY="+ready)
			var combined bytes.Buffer
			cmd.Stdout, cmd.Stderr = &combined, &combined
			if err := cmd.Start(); err != nil {
				t.Fatal(err)
			}
			defer func() {
				if cmd.ProcessState == nil {
					_ = cmd.Process.Kill()
					_ = cmd.Wait()
				}
			}()
			deadline := time.Now().Add(2 * time.Second)
			for {
				if _, err := os.Stat(ready); err == nil {
					break
				}
				if time.Now().After(deadline) {
					t.Fatalf("child never became ready: %q", combined.String())
				}
				time.Sleep(10 * time.Millisecond)
			}
			if _, err := os.Stat(trace); !os.IsNotExist(err) {
				t.Fatalf("handler ran before SIGINT: trace stat %v", err)
			}
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				t.Fatal(err)
			}
			err := cmd.Wait()
			if ctx.Err() != nil {
				t.Fatalf("child timed out: %v", ctx.Err())
			}
			if status := exitCode(err); status != tc.wantStatus {
				t.Fatalf("status %d, want %d; output %q", status, tc.wantStatus, combined.String())
			}
			gotTrace, err := os.ReadFile(trace)
			if err != nil || string(gotTrace) != tc.wantTrace {
				t.Fatalf("trace %q, err %v, want %q", gotTrace, err, tc.wantTrace)
			}
			if tc.wantOutput != "" && combined.String() != tc.wantOutput {
				t.Fatalf("output %q, want %q", combined.String(), tc.wantOutput)
			}
		})
	}
	invalid := filepath.Join(t.TempDir(), "invalid.sh")
	mustWriteFile(t, invalid, []byte(expandedGeneratedCtrlC43(t, `printf 'unterminated`)+"\n"))
	if _, err := runCommand(".", nil, "sh", "-n", invalid); err == nil {
		t.Fatal("malformed inserted source unexpectedly passed syntax check")
	}
}

func expandedGeneratedExit44(t *testing.T, action string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration44MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentExitBody) {
		t.Fatalf("generated EXIT body = %#v", snippets[iteration44MigrationName].Body)
	}
	lines := make([]string, len(body))
	for index, element := range body {
		line, ok := element.(string)
		if !ok {
			t.Fatalf("EXIT line %d is not a string", index)
		}
		lines[index] = line
	}
	joined := strings.Join(lines, "\n")
	placeholder := `${1:# i.e. clean-up code here}`
	if strings.Count(joined, placeholder) != 1 ||
		!strings.Contains(joined, "on_exit() {\n") ||
		!strings.Contains(joined, "\n\texit 0 # Exit gracefully.") ||
		!strings.HasSuffix(joined, "trap on_exit EXIT") {
		t.Fatalf("unexpected generated EXIT syntax: %q", joined)
	}
	return strings.Replace(joined, placeholder, action, 1)
}

func TestExit44GeneratedExpansionAndStatus(t *testing.T) {
	for _, shell := range []string{"sh", "bash", "zsh"} {
		requireCommand(t, shell)
		path := filepath.Join(t.TempDir(), "syntax.sh")
		mustWriteFile(t, path, []byte(expandedGeneratedExit44(t, "# i.e. clean-up code here")+"\n"))
		run(t, ".", shell, "-n", path)
	}
	defaultPath := filepath.Join(t.TempDir(), "default.sh")
	mustWriteFile(t, defaultPath, []byte("#!/bin/sh\n"+expandedGeneratedExit44(t, "# i.e. clean-up code here")+"\n"))
	output, err := runCommand(".", nil, "shellcheck", "-s", "sh", defaultPath)
	if err != nil {
		t.Fatalf("ShellCheck on exact generated EXIT expansion: %v: %s", err, output)
	}

	type exitCase struct {
		name, shell, before, after, action, wantTrace string
		wantStatus                                    int
		fakeTput, captureEntry                        bool
	}
	cases := []exitCase{
		{name: "end of script zero", shell: "sh", after: ":", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "ENTRY:0\nACTION\n", captureEntry: true},
		{name: "end of script failure", shell: "sh", after: "false", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "ENTRY:1\nACTION\n", captureEntry: true},
		{name: "explicit exit zero", shell: "sh", after: "exit 0", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "ENTRY:0\nACTION\n", captureEntry: true},
		{name: "explicit exit one", shell: "sh", after: "exit 1", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "ENTRY:1\nACTION\n", captureEntry: true},
		{name: "explicit exit forty-two", shell: "sh", after: "exit 42", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "ENTRY:42\nACTION\n", captureEntry: true},
		{name: "set-e controlled failure", shell: "sh", after: "set -e\nfalse", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "ENTRY:1\nACTION\n", captureEntry: true},
		{name: "prior exit trap replaced", shell: "sh", before: `trap 'printf "OLD\n" >> "$TRACE"' EXIT`, after: ":", action: `printf 'NEW\n' >> "$TRACE"`, wantTrace: "TPUT\nNEW\n", fakeTput: true},
		{name: "late value substitution multiple commands", shell: "sh", before: "value=before", after: "value=after\n:", action: `printf '%s\n' "$(printf 'SUB\n' >> "$TRACE"; printf '%s' "$value")" >> "$TRACE"; printf 'SECOND\n' >> "$TRACE"`, wantTrace: "TPUT\nSUB\nafter\nSECOND\n", fakeTput: true},
		{name: "quotes and arithmetic", shell: "sh", after: ":", action: `printf '%s %s\n' 'two words' "$((2+3))" >> "$TRACE"`, wantTrace: "TPUT\ntwo words 5\n", fakeTput: true},
		{name: "empty placeholder", shell: "sh", after: "exit 42", action: "", wantTrace: "TPUT\n", fakeTput: true},
		{name: "failed action followed by success", shell: "sh", after: "exit 42", action: `false; printf 'LATER\n' >> "$TRACE"`, wantTrace: "TPUT\nLATER\n", fakeTput: true},
		{name: "handler explicit exit zero", shell: "sh", after: "exit 42", action: `exit 0`, wantTrace: "TPUT\n", fakeTput: true},
		{name: "handler explicit exit forty-two", shell: "sh", after: "exit 1", action: `exit 42`, wantTrace: "TPUT\n", wantStatus: 42, fakeTput: true},
		{name: "handler returns successfully", shell: "sh", after: "exit 42", action: `return 0`, wantTrace: "ENTRY:42\n", wantStatus: 42, captureEntry: true},
		{name: "handler returns failure", shell: "sh", after: "exit 42", action: `return 7`, wantTrace: "ENTRY:42\n", wantStatus: 42, captureEntry: true},
		{name: "subshell does not inherit exit trap", shell: "sh", after: "(exit 42)\nexit 1", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "TPUT\nACTION\n", fakeTput: true},
		{name: "unavailable action command", shell: "sh", after: ":", action: `missing_shellman_exit_command; printf 'LATER\n' >> "$TRACE"`, wantTrace: "TPUT\nLATER\n", fakeTput: true},
		{name: "unavailable tput", shell: "sh", after: "exit 42", action: `printf 'TPUT_STATUS:%s\n' "$?" >> "$TRACE"`, wantTrace: "TPUT_STATUS:127\n"},
		{name: "bash runtime", shell: "bash", after: "exit 42", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "TPUT\nACTION\n", fakeTput: true},
		{name: "zsh runtime", shell: "zsh", after: "exit 42", action: `printf 'ACTION\n' >> "$TRACE"`, wantTrace: "TPUT\nACTION\n", fakeTput: true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			bin := filepath.Join(root, "bin")
			mustMkdirAll(t, bin)
			if tc.fakeTput {
				mustExecutable(t, filepath.Join(bin, "tput"), "#!/bin/sh\n[ \"$#\" -eq 1 ] && [ \"$1\" = cnorm ] || exit 42\nprintf 'TPUT\\n' >> \"$TRACE\"\n")
			}
			before := tc.before
			if tc.captureEntry {
				before += "\ntput() { entry_status=$?; printf 'ENTRY:%s\\n' \"$entry_status\" >> \"$TRACE\"; return 0; }"
			}
			trace := filepath.Join(root, "trace")
			path := filepath.Join(root, "child.sh")
			mustWriteFile(t, path, []byte("#!/bin/sh\n"+before+"\n"+expandedGeneratedExit44(t, tc.action)+"\n"+tc.after+"\n"))
			run(t, ".", tc.shell, "-n", path)
			ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
			defer cancel()
			cmd := exec.CommandContext(ctx, tc.shell, path)
			cmd.Env = append(os.Environ(), "PATH="+bin, "TRACE="+trace)
			output, err := cmd.CombinedOutput()
			if ctx.Err() != nil {
				t.Fatalf("child timeout: %v", ctx.Err())
			}
			if status := exitCode(err); status != tc.wantStatus {
				t.Fatalf("exit status %d, want %d: %q", status, tc.wantStatus, output)
			}
			gotTrace, err := os.ReadFile(trace)
			if err != nil || string(gotTrace) != tc.wantTrace {
				t.Fatalf("trace %q, err %v, want %q; output %q", gotTrace, err, tc.wantTrace, output)
			}
		})
	}
	t.Run("caught SIGINT exits and then runs EXIT trap", func(t *testing.T) {
		root := t.TempDir()
		bin := filepath.Join(root, "bin")
		mustMkdirAll(t, bin)
		mustExecutable(t, filepath.Join(bin, "tput"), "#!/bin/sh\nprintf 'TPUT\\n' >> \"$TRACE\"\n")
		trace, ready := filepath.Join(root, "trace"), filepath.Join(root, "ready")
		path := filepath.Join(root, "child.sh")
		body := expandedGeneratedExit44(t, `printf 'EXIT\n' >> "$TRACE"`)
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+body+"\ntrap 'exit 42' INT\nprintf 'ready\\n' > \"$READY\"\nwhile :; do :; done\n"))
		run(t, ".", "sh", "-n", path)
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "sh", path)
		cmd.Env = append(os.Environ(), "PATH="+bin, "TRACE="+trace, "READY="+ready)
		var combined bytes.Buffer
		cmd.Stdout, cmd.Stderr = &combined, &combined
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if cmd.ProcessState == nil {
				_ = cmd.Process.Kill()
				_ = cmd.Wait()
			}
		}()
		deadline := time.Now().Add(2 * time.Second)
		for {
			if _, err := os.Stat(ready); err == nil {
				break
			}
			if time.Now().After(deadline) {
				t.Fatalf("child not ready: %q", combined.String())
			}
			time.Sleep(10 * time.Millisecond)
		}
		if _, err := os.Stat(trace); !os.IsNotExist(err) {
			t.Fatalf("EXIT trap ran before SIGINT: trace stat %v", err)
		}
		if err := cmd.Process.Signal(os.Interrupt); err != nil {
			t.Fatal(err)
		}
		err := cmd.Wait()
		if ctx.Err() != nil || exitCode(err) != 0 {
			t.Fatalf("status %d, timeout %v, output %q", exitCode(err), ctx.Err(), combined.String())
		}
		got, err := os.ReadFile(trace)
		if err != nil || string(got) != "TPUT\nEXIT\n" {
			t.Fatalf("trace %q, err %v, output %q", got, err, combined.String())
		}
	})
	t.Run("successful exec replacement does not exit shell", func(t *testing.T) {
		root := t.TempDir()
		trace := filepath.Join(root, "trace")
		path := filepath.Join(root, "child.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedExit44(t, `printf 'EXIT\n' >> "$TRACE"`)+"\nexec /bin/sh -c 'exit 37'\n"))
		run(t, ".", "sh", "-n", path)
		ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "sh", path)
		cmd.Env = append(os.Environ(), "TRACE="+trace)
		output, err := cmd.CombinedOutput()
		if ctx.Err() != nil || exitCode(err) != 37 {
			t.Fatalf("status %d, timeout %v, output %q", exitCode(err), ctx.Err(), output)
		}
		if _, err := os.Stat(trace); !os.IsNotExist(err) {
			t.Fatalf("EXIT handler ran after exec replacement: trace stat %v", err)
		}
	})
	invalid := filepath.Join(t.TempDir(), "invalid.sh")
	mustWriteFile(t, invalid, []byte(expandedGeneratedExit44(t, `printf 'unterminated`)+"\n"))
	if _, err := runCommand(".", nil, "sh", "-n", invalid); err == nil {
		t.Fatal("malformed inserted shell source unexpectedly passed syntax check")
	}
}

func expandedGeneratedDirectoryIterate45(t *testing.T, choice, action string) string {
	t.Helper()
	generated, _, err := generate(".")
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(generated, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration45MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentDirectoriesIterateBody) {
		t.Fatalf("generated directory body = %#v", snippets[iteration45MigrationName].Body)
	}
	lines := make([]string, len(body))
	for index, element := range body {
		line, ok := element.(string)
		if !ok {
			t.Fatalf("directory body line %d is not a string", index)
		}
		lines[index] = line
	}
	joined := strings.Join(lines, "\n")
	pathPlaceholder := `${1| /'path/to/directory'/,"${pathToDirectory}"|}`
	actionPlaceholder := `${2:printf '%s\n' "\${directory\}"}`
	if strings.Count(joined, pathPlaceholder) != 1 || strings.Count(joined, actionPlaceholder) != 1 ||
		!strings.HasSuffix(joined, "done\n") {
		t.Fatalf("unexpected generated directory placeholders: %q", joined)
	}
	if choice != ` /'path/to/directory'/` && choice != `"${pathToDirectory}"` {
		t.Fatalf("invalid path choice %q", choice)
	}
	joined = strings.Replace(joined, pathPlaceholder, choice, 1)
	joined = strings.Replace(joined, actionPlaceholder, action, 1)
	return strings.ReplaceAll(joined, `\${directory}`, `${directory}`)
}

func TestDirectoriesIterate45GeneratedExpansionAndPaths(t *testing.T) {
	requireCommand(t, "sh")
	defaultAction := `printf '%s\n' "${directory}"`
	variableChoice := `"${pathToDirectory}"`
	defaultBody := expandedGeneratedDirectoryIterate45(t, variableChoice, defaultAction)
	wantDefault := "# Make sure path ends with /\nfor directory in \"${pathToDirectory}\"*; do\n\tif [ -d \"${directory}\" ] && [ ! -L \"${directory}\" ]; then\n\t\tprintf '%s\\n' \"${directory}\"\n\tfi\ndone\n"
	if defaultBody != wantDefault {
		t.Fatalf("exact default expansion = %q, want %q", defaultBody, wantDefault)
	}
	if got := expandedGeneratedDirectoryIterate45(t, ` /'path/to/directory'/`, defaultAction); !strings.Contains(got, "for directory in  /'path/to/directory'/*; do") {
		t.Fatalf("unexpected literal path choice expansion: %q", got)
	}
	defaultPath := filepath.Join(t.TempDir(), "default.sh")
	mustWriteFile(t, defaultPath, []byte("#!/bin/sh\npathToDirectory=/nonexistent/\n"+defaultBody))
	for _, shell := range []string{"sh", "bash", "zsh"} {
		requireCommand(t, shell)
		run(t, ".", shell, "-n", defaultPath)
	}
	output, err := runCommand(".", nil, "shellcheck", "-s", "sh", defaultPath)
	if err != nil {
		t.Fatalf("ShellCheck on generated directory expansion: %v: %s", err, output)
	}

	root := t.TempDir()
	base := filepath.Join(root, "parent [x]* $;()")
	mustMkdirAll(t, base)
	names := []string{"plain", "with space", "with\ttab", "wild*?", "-leading", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak"}
	for _, name := range names {
		mustMkdirAll(t, filepath.Join(base, name))
	}
	mustMkdirAll(t, filepath.Join(base, "plain", "nested"))
	mustMkdirAll(t, filepath.Join(base, ".hidden"))
	mustWriteFile(t, filepath.Join(base, "ordinary-file"), []byte("file"))
	if err := os.Symlink(filepath.Join(base, "plain"), filepath.Join(base, "link-to-directory")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(base, "missing"), filepath.Join(base, "broken-link")); err != nil {
		t.Fatal(err)
	}
	trace := filepath.Join(root, "trace")
	path := filepath.Join(root, "iterate.sh")
	action := `printf '%s\000' "$directory" >> "$TRACE"`
	mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedDirectoryIterate45(t, variableChoice, action)))
	run(t, ".", "sh", "-n", path)
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "sh", path)
	cmd.Env = append(os.Environ(), "pathToDirectory="+base+"/", "TRACE="+trace, "LC_ALL=C")
	output, err = cmd.CombinedOutput()
	if err != nil || ctx.Err() != nil {
		t.Fatalf("adversarial directory run: %v, timeout %v, output %q", err, ctx.Err(), output)
	}
	gotBytes, err := os.ReadFile(trace)
	if err != nil {
		t.Fatal(err)
	}
	if len(gotBytes) == 0 || gotBytes[len(gotBytes)-1] != 0 {
		t.Fatalf("missing final NUL in trace %q", gotBytes)
	}
	got := strings.Split(string(gotBytes[:len(gotBytes)-1]), "\x00")
	want := make([]string, len(names))
	for index, name := range names {
		want[index] = filepath.Join(base, name)
	}
	sort.Strings(got)
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("pathnames changed or wrong membership: got %#v, want %#v", got, want)
	}
	for _, shell := range []string{"sh", "zsh"} {
		t.Run("default-print-"+shell, func(t *testing.T) {
			root := t.TempDir()
			mustMkdirAll(t, filepath.Join(root, `back\slash`))
			path := filepath.Join(root, "print.sh")
			mustWriteFile(t, path, []byte("#!/bin/sh\n"+defaultBody))
			output, err := runCommand(".", append(os.Environ(), "pathToDirectory="+root+"/"), shell, path)
			if err != nil || string(output) != filepath.Join(root, `back\slash`)+"\n" {
				t.Fatalf("default path output = %q, err %v", output, err)
			}
		})
	}
	invalid := filepath.Join(t.TempDir(), "invalid.sh")
	mustWriteFile(t, invalid, []byte(expandedGeneratedDirectoryIterate45(t, variableChoice, `printf 'unterminated`)))
	if _, err := runCommand(".", nil, "sh", "-n", invalid); err == nil {
		t.Fatal("malformed inserted action unexpectedly passed syntax check")
	}
}

func TestDirectoriesIterate45ActionAndFailureStatus(t *testing.T) {
	choice := `"${pathToDirectory}"`
	runCase := func(t *testing.T, entries []string, action, before, after string) (string, int) {
		t.Helper()
		root := t.TempDir()
		base := filepath.Join(root, "source")
		mustMkdirAll(t, base)
		for _, name := range entries {
			if strings.HasPrefix(name, "file:") {
				mustWriteFile(t, filepath.Join(base, strings.TrimPrefix(name, "file:")), []byte("regular"))
			} else {
				mustMkdirAll(t, filepath.Join(base, name))
			}
		}
		path := filepath.Join(root, "iterate.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+before+"\n"+expandedGeneratedDirectoryIterate45(t, choice, action)+after+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "pathToDirectory="+base+"/", "LC_ALL=C"), "sh", path)
		return string(output), exitCode(err)
	}
	t.Run("zero matches and hidden only", func(t *testing.T) {
		output, status := runCase(t, []string{".hidden", "file:visible"}, `printf 'ACTION\n'`, "", "exit $? ")
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("one directory", func(t *testing.T) {
		output, status := runCase(t, []string{"one"}, `printf '%s\n' "${directory##*/}"`, "", "exit $? ")
		if status != 0 || output != "one\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("assignments remain in current shell", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "z"}, `count=$((count+1))`, "count=0", `printf 'COUNT:%s\n' "$count"`)
		if status != 0 || output != "COUNT:2\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("earlier action failure masked", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "z"}, `case ${directory##*/} in a) false;; *) :;; esac`, "", "exit $? ")
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("last action failure controls loop status", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "z"}, `case ${directory##*/} in z) false;; *) :;; esac`, "", "exit $? ")
		if status != 1 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("later skipped file masks action failure", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "file:z"}, `false`, "", "exit $? ")
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("action command not found", func(t *testing.T) {
		output, status := runCase(t, []string{"one"}, `missing_shellman_iteration45_command`, "", "exit $? ")
		if status != 127 || !strings.Contains(output, "missing_shellman_iteration45_command") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("break leaves loop", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "z"}, `printf '%s\n' "${directory##*/}"; break`, "", "exit $? ")
		if status != 0 || output != "a\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("continue skips rest of action", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "z"}, `if [ "${directory##*/}" = a ]; then continue; fi; printf '%s\n' "${directory##*/}"`, "", "exit $? ")
		if status != 0 || output != "z\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("return leaves enclosing function", func(t *testing.T) {
		root := t.TempDir()
		base := filepath.Join(root, "source")
		mustMkdirAll(t, filepath.Join(base, "a"))
		mustMkdirAll(t, filepath.Join(base, "z"))
		path := filepath.Join(root, "return.sh")
		body := expandedGeneratedDirectoryIterate45(t, choice, `return 7`)
		mustWriteFile(t, path, []byte("#!/bin/sh\nwalk() {\n"+body+"}\nwalk\nexit $?\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "pathToDirectory="+base+"/"), "sh", path)
		if status := exitCode(err); status != 7 || len(output) != 0 {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("exit leaves shell", func(t *testing.T) {
		output, status := runCase(t, []string{"a", "z"}, `exit 42`, "", `printf 'AFTER\n'`)
		if status != 42 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("nonexistent starting path", func(t *testing.T) {
		root := t.TempDir()
		path := filepath.Join(root, "missing.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedDirectoryIterate45(t, choice, `printf 'ACTION\n'`)))
		output, err := runCommand(".", append(os.Environ(), "pathToDirectory="+filepath.Join(root, "absent")+"/"), "sh", path)
		if status := exitCode(err); status != 0 || len(output) != 0 {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("relative starting path with leading hyphen", func(t *testing.T) {
		root := t.TempDir()
		mustMkdirAll(t, filepath.Join(root, "-root", "inside"))
		path := filepath.Join(root, "relative.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedDirectoryIterate45(t, choice, `printf '%s\n' "$directory"`)))
		output, err := runCommand(root, append(os.Environ(), "pathToDirectory=-root/"), "sh", path)
		if status := exitCode(err); status != 0 || string(output) != "-root/inside\n" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
	t.Run("unreadable starting directory", func(t *testing.T) {
		if os.Geteuid() == 0 {
			t.Skip("permission test is not meaningful as root")
		}
		root := t.TempDir()
		base := filepath.Join(root, "unreadable")
		mustMkdirAll(t, filepath.Join(base, "inside"))
		if err := os.Chmod(base, 0); err != nil {
			t.Fatal(err)
		}
		defer os.Chmod(base, 0o755)
		if _, err := os.ReadDir(base); err == nil {
			t.Skip("current privileges still permit reading mode-000 directory")
		}
		path := filepath.Join(root, "unreadable.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+expandedGeneratedDirectoryIterate45(t, choice, `printf 'ACTION\n'`)))
		output, err := runCommand(".", append(os.Environ(), "pathToDirectory="+base+"/"), "sh", path)
		if status := exitCode(err); status != 0 || len(output) != 0 {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
}

func TestDateNowShortGeneratedExpansionAndStatus(t *testing.T) {
	requireCommand(t, "sh")
	if got, want := expandedGeneratedDateNowShort(t, "date_result"), "date_result=$(date '+%Y/%m/%d') # format: yyyy/mm/dd\n"; got != want {
		t.Fatalf("expanded generated body = %q, want %q", got, want)
	}

	runCase := func(t *testing.T, stub, setup, destination, trailer string, noDate bool) (string, int) {
		t.Helper()
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		if !noDate {
			mustExecutable(t, filepath.Join(fakeBin, "date"), "#!/bin/sh\n"+stub)
		}
		path := filepath.Join(root, "snippet.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n"+setup+"\n"+expandedGeneratedDateNowShort(t, destination)+"snippet_status=$?\n"+trailer+"\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "PATH="+fakeBin), "sh", path)
		return string(output), exitCode(err)
	}

	t.Run("exact operand and newline stripping", func(t *testing.T) {
		stub := "[ \"$#\" -eq 1 ] && [ \"$1\" = '+%Y/%m/%d' ] || exit 42\nprintf '2034/05/06\\n\\n'\n"
		output, status := runCase(t, stub, "", "date_result", `[ "$snippet_status" -eq 0 ] && [ "$date_result" = '2034/05/06' ]`, false)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("leading and internal newlines survive", func(t *testing.T) {
		stub := "printf '\\n2034/05\\n06\\n'\n"
		output, status := runCase(t, stub, "", "date_result", `[ "$snippet_status" -eq 0 ] && [ "$date_result" = "$(printf '\n2034/05\n06')" ]`, false)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	for _, failureStatus := range []int{1, 42, 127} {
		t.Run(fmt.Sprintf("utility failure %d", failureStatus), func(t *testing.T) {
			stub := fmt.Sprintf("printf 'partial\\n'; printf 'date error\\n' >&2; exit %d\n", failureStatus)
			output, status := runCase(t, stub, "", "date_result", `[ "$date_result" = partial ] && exit "$snippet_status"`, false)
			if status != failureStatus || !strings.Contains(output, "date error") {
				t.Fatalf("status=%d output=%q", status, output)
			}
		})
	}

	t.Run("actual command not found", func(t *testing.T) {
		output, status := runCase(t, "", "", "date_result", `exit "$snippet_status"`, true)
		if status != 127 || !strings.Contains(output, "date") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("readonly destination", func(t *testing.T) {
		output, status := runCase(t, "printf '2034/05/06\\n'\n", "readonly date_result", "date_result", `exit "$snippet_status"`, false)
		if status == 0 || !strings.Contains(output, "date_result") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("invalid destination", func(t *testing.T) {
		output, status := runCase(t, "printf '2034/05/06\\n'\n", "", "1bad", `exit "$snippet_status"`, false)
		if status == 0 || !strings.Contains(output, "1bad") {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("later command masks status", func(t *testing.T) {
		output, status := runCase(t, "exit 42\n", "", "date_result", `:; [ "$snippet_status" -eq 42 ]`, false)
		if status != 0 || output != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})

	t.Run("pipeline masks left command status", func(t *testing.T) {
		root := t.TempDir()
		fakeBin := filepath.Join(root, "bin")
		mustMkdirAll(t, fakeBin)
		mustExecutable(t, filepath.Join(fakeBin, "date"), "#!/bin/sh\nexit 42\n")
		path := filepath.Join(root, "pipeline.sh")
		mustWriteFile(t, path, []byte("#!/bin/sh\n{\n"+expandedGeneratedDateNowShort(t, "date_result")+"} | /bin/cat\n"))
		run(t, ".", "sh", "-n", path)
		output, err := runCommand(".", append(os.Environ(), "PATH="+fakeBin), "sh", path)
		status := exitCode(err)
		if status != 0 || string(output) != "" {
			t.Fatalf("status=%d output=%q", status, output)
		}
	})
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

func expandedGeneratedDirectoryCreateNested46(t *testing.T, choice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration46MigrationName].Body.(string)
	if !ok || body != currentDirectoryCreateNestedBody {
		t.Fatalf("generated entry 46 body = %#v", snippets[iteration46MigrationName].Body)
	}
	placeholder := `${1|"parent dir/child dir","${pathToCreate}"|}`
	if strings.Count(body, placeholder) != 1 {
		t.Fatalf("unexpected entry 46 placeholder: %q", body)
	}
	return strings.Replace(body, placeholder, choice, 1)
}

func TestDirectoryCreateNested46GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	defaultBody := expandedGeneratedDirectoryCreateNested46(t, `"parent dir/child dir"`)
	if want := "directory_path=\"parent dir/child dir\"\ncase $directory_path in\n  -*) directory_path=./$directory_path ;;\nesac\nmkdir -p \"$directory_path\"\n"; defaultBody != want {
		t.Fatalf("default expansion = %q, want %q", defaultBody, want)
	}
	root := t.TempDir()
	defaultScript := filepath.Join(root, "default.sh")
	mustWriteFile(t, defaultScript, []byte("#!/bin/sh\n"+defaultBody))
	run(t, ".", "sh", "-n", defaultScript)
	if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", defaultScript); err != nil {
		t.Fatalf("ShellCheck: %v: %s", err, output)
	}
	if output, err := runCommand(root, nil, "sh", defaultScript); err != nil || len(output) != 0 {
		t.Fatalf("default creation: %v, %q", err, output)
	}
	if info, err := os.Stat(filepath.Join(root, "parent dir", "child dir")); err != nil || !info.IsDir() {
		t.Fatalf("default target: %v", err)
	}
	variableBody := expandedGeneratedDirectoryCreateNested46(t, `"${pathToCreate}"`)
	variableScript := filepath.Join(root, "variable.sh")
	mustWriteFile(t, variableScript, []byte("#!/bin/sh\n# shellcheck disable=SC2154\n"+variableBody))
	run(t, ".", "sh", "-n", variableScript)
	if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", variableScript); err != nil {
		t.Fatalf("ShellCheck variable choice: %v: %s", err, output)
	}
	runCase := func(path string) (string, int) {
		t.Helper()
		output, err := runCommand(root, append(os.Environ(), "pathToCreate="+path), "sh", variableScript)
		return string(output), exitCode(err)
	}
	for _, name := range []string{"space name", "tab\tname", "wild*?name", "-leading", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak"} {
		path := name + "/nested"
		output, status := runCase(path)
		if status != 0 || output != "" {
			t.Fatalf("create %q: status %d output %q", path, status, output)
		}
		if info, err := os.Stat(filepath.Join(root, name, "nested")); err != nil || !info.IsDir() {
			t.Fatalf("path %q: %v", path, err)
		}
	}
	for _, path := range []string{"-leading/nested", "space name/nested", "./relative/../relative/child", "trailing/child/", "."} {
		if output, status := runCase(path); status != 0 || output != "" {
			t.Fatalf("repeat or normalized %q: %d %q", path, status, output)
		}
	}
	abs := filepath.Join(root, "absolute", "child")
	if output, status := runCase(abs); status != 0 || output != "" {
		t.Fatalf("absolute path: %d %q", status, output)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase("/"); status != 0 || output != "" {
		t.Fatalf("root target: %d %q", status, output)
	}
	mustWriteFile(t, filepath.Join(root, "regular"), []byte("file"))
	for _, path := range []string{"regular", "regular/child", ""} {
		if output, status := runCase(path); status == 0 || output == "" {
			t.Fatalf("failure %q: %d %q", path, status, output)
		}
	}
	if err := os.Symlink("space name", filepath.Join(root, "dirlink")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("regular", filepath.Join(root, "filelink")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("absent", filepath.Join(root, "broken")); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase("dirlink"); status != 0 || output != "" {
		t.Fatalf("directory symlink: %d %q", status, output)
	}
	if output, status := runCase("dirlink/new"); status != 0 || output != "" {
		t.Fatalf("intermediate directory symlink: %d %q", status, output)
	}
	if _, err := os.Stat(filepath.Join(root, "space name", "new")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"filelink", "filelink/new", "broken", "broken/new"} {
		if output, status := runCase(path); status == 0 || output == "" {
			t.Fatalf("symlink failure %q: %d %q", path, status, output)
		}
	}
	if output, status := runCase("partial/../regular/child"); status == 0 || output == "" {
		t.Fatalf("partial failure: %d %q", status, output)
	}
	if _, err := os.Stat(filepath.Join(root, "partial")); err != nil {
		t.Fatal(err)
	}
	modeScript := filepath.Join(root, "mode.sh")
	mustWriteFile(t, modeScript, []byte("#!/bin/sh\numask 027\n"+variableBody))
	if output, err := runCommand(root, append(os.Environ(), "pathToCreate=mode-parent/child"), "sh", modeScript); err != nil || len(output) != 0 {
		t.Fatalf("mode creation: %v %q", err, output)
	}
	for _, path := range []string{"mode-parent", "mode-parent/child"} {
		info, err := os.Stat(filepath.Join(root, path))
		if err != nil || info.Mode().Perm() != 0o750 {
			t.Fatalf("mode for %q: %v, %v", path, info, err)
		}
	}
	if output, status := runCase("mode-parent/child"); status != 0 || output != "" {
		t.Fatalf("existing target status: %d %q", status, output)
	}
	if info, err := os.Stat(filepath.Join(root, "mode-parent/child")); err != nil || info.Mode().Perm() != 0o750 {
		t.Fatalf("existing mode changed: %v, %v", info, err)
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "blocked")
		mustMkdirAll(t, blocked)
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runCase("blocked/child")
		if err := os.Chmod(blocked, 0o700); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("permission failure: %d %q", status, output)
		}
	}
	missingOutput, missingErr := runCommand(root, []string{"PATH=/nonexistent", "pathToCreate=missing-utility"}, "/bin/sh", variableScript)
	if exitCode(missingErr) != 127 || !strings.Contains(string(missingOutput), "mkdir") {
		t.Fatalf("missing utility: %v %q", missingErr, missingOutput)
	}
	masked, err := runCommand(root, append(os.Environ(), "pathToCreate=regular/child"), "sh", "-c", variableBody+"true")
	if err != nil || len(masked) == 0 {
		t.Fatalf("later command masking: %v %q", err, masked)
	}
}

func expandedGeneratedDirectoryCreate47(t *testing.T, replacement string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration47MigrationName].Body.(string)
	if !ok || body != currentDirectoryCreateBody {
		t.Fatalf("generated entry-47 body = %#v", snippets[iteration47MigrationName].Body)
	}
	placeholder := `${1:directory name}`
	if strings.Count(body, placeholder) != 1 {
		t.Fatalf("entry-47 placeholder = %q", body)
	}
	return strings.Replace(body, placeholder, replacement, 1)
}

func TestDirectoryCreate47GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	root := t.TempDir()
	defaultBody := expandedGeneratedDirectoryCreate47(t, "directory name")
	wantDefault := "directory_path=\"directory name\"\ncase $directory_path in\n  -*) directory_path=./$directory_path ;;\nesac\nmkdir \"$directory_path\"\n"
	if defaultBody != wantDefault {
		t.Fatalf("default expansion = %q, want %q", defaultBody, wantDefault)
	}
	defaultScript := filepath.Join(root, "default.sh")
	mustWriteFile(t, defaultScript, []byte("#!/bin/sh\n"+defaultBody))
	run(t, ".", "sh", "-n", defaultScript)
	if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", defaultScript); err != nil {
		t.Fatalf("default lint: %v %q", err, output)
	}
	if output, err := runCommand(root, nil, "sh", defaultScript); err != nil || len(output) != 0 {
		t.Fatalf("default create: %v %q", err, output)
	}
	if info, err := os.Stat(filepath.Join(root, "directory name")); err != nil || !info.IsDir() {
		t.Fatalf("default target: %v", err)
	}

	variableBody := expandedGeneratedDirectoryCreate47(t, `$directoryName`)
	variableScript := filepath.Join(root, "variable.sh")
	mustWriteFile(t, variableScript, []byte("#!/bin/sh\n# shellcheck disable=SC2154\n"+variableBody))
	run(t, ".", "sh", "-n", variableScript)
	if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", variableScript); err != nil {
		t.Fatalf("variable lint: %v %q", err, output)
	}
	runCase := func(path string) (string, int) {
		t.Helper()
		output, err := runCommand(root, append(os.Environ(), "directoryName="+path), "sh", variableScript)
		return string(output), exitCode(err)
	}
	for _, name := range []string{"with space", "with\ttab", "wild*?card", "-leading", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak"} {
		output, status := runCase(name)
		if status != 0 || output != "" {
			t.Fatalf("create %q: %d %q", name, status, output)
		}
		if info, err := os.Stat(filepath.Join(root, name)); err != nil || !info.IsDir() {
			t.Fatalf("path %q: %v", name, err)
		}
	}
	for _, path := range []string{"directory name", "-leading", "with space", ".", "..", "/", ""} {
		output, status := runCase(path)
		if status == 0 || output == "" {
			t.Fatalf("expected failure %q: %d %q", path, status, output)
		}
	}
	for _, path := range []string{"missing/child", "missing/child/", "regular/child"} {
		output, status := runCase(path)
		if status == 0 || output == "" {
			t.Fatalf("missing/non-directory parent %q: %d %q", path, status, output)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "missing")); !os.IsNotExist(err) {
		t.Fatalf("missing parent was created: %v", err)
	}
	abs := filepath.Join(root, "absolute-target")
	if output, status := runCase(abs); status != 0 || output != "" {
		t.Fatalf("absolute: %d %q", status, output)
	}
	if _, err := os.Stat(abs); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase("trailing/"); status != 0 || output != "" {
		t.Fatalf("trailing slash: %d %q", status, output)
	}
	if output, status := runCase("./relative"); status != 0 || output != "" {
		t.Fatalf("relative: %d %q", status, output)
	}
	if output, status := runCase("relative/../dotdot"); status != 0 || output != "" {
		t.Fatalf("dotdot: %d %q", status, output)
	}
	mustWriteFile(t, filepath.Join(root, "regular"), []byte("file"))
	if output, status := runCase("regular"); status == 0 || output == "" {
		t.Fatalf("existing file: %d %q", status, output)
	}
	if output, status := runCase("regular/child"); status == 0 || output == "" {
		t.Fatalf("regular-file parent: %d %q", status, output)
	}
	if err := os.Symlink("with space", filepath.Join(root, "dirlink")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("regular", filepath.Join(root, "filelink")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("absent", filepath.Join(root, "broken")); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{"dirlink", "filelink", "broken", "filelink/child", "broken/child"} {
		if output, status := runCase(path); status == 0 || output == "" {
			t.Fatalf("symlink failure %q: %d %q", path, status, output)
		}
	}
	if output, status := runCase("dirlink/child"); status != 0 || output != "" {
		t.Fatalf("directory symlink parent: %d %q", status, output)
	}
	if _, err := os.Stat(filepath.Join(root, "with space", "child")); err != nil {
		t.Fatal(err)
	}
	modeScript := filepath.Join(root, "mode.sh")
	mustWriteFile(t, modeScript, []byte("#!/bin/sh\numask 027\n"+variableBody+"umask\n"))
	output, err := runCommand(root, append(os.Environ(), "directoryName=mode-target"), "sh", modeScript)
	if err != nil || !strings.Contains(string(output), "0027") {
		t.Fatalf("umask status/value: %v %q", err, output)
	}
	if info, err := os.Stat(filepath.Join(root, "mode-target")); err != nil || info.Mode().Perm() != 0o750 {
		t.Fatalf("mode: %v %v", info, err)
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "blocked")
		mustMkdirAll(t, blocked)
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runCase("blocked/child")
		if err := os.Chmod(blocked, 0o700); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("permission denial: %d %q", status, output)
		}
	}
	missingUtility, err := runCommand(root, []string{"PATH=/nonexistent", "directoryName=not-created"}, "/bin/sh", variableScript)
	if exitCode(err) != 127 || !strings.Contains(string(missingUtility), "mkdir") {
		t.Fatalf("missing utility: %v %q", err, missingUtility)
	}
	masked, err := runCommand(root, append(os.Environ(), "directoryName=directory name"), "sh", "-c", variableBody+"true")
	if err != nil || len(masked) == 0 {
		t.Fatalf("later command masks failure: %v %q", err, masked)
	}
}

func expandedGeneratedDirectoryDelete48(t *testing.T, choice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration48MigrationName].Body.(string)
	if !ok || body != currentDirectoryDeleteNestedBody {
		t.Fatalf("generated entry 48: %#v", snippets[iteration48MigrationName].Body)
	}
	placeholder := `${1|/path/to/directory,${pathToDirectory}|}`
	if strings.Count(body, placeholder) != 1 {
		t.Fatalf("placeholder count in %q", body)
	}
	return strings.Replace(body, placeholder, choice, 1)
}

func TestDirectoryDelete48GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	defaultBody := expandedGeneratedDirectoryDelete48(t, "/path/to/directory")
	variableBody := expandedGeneratedDirectoryDelete48(t, `${pathToDirectory}`)
	for _, body := range []string{defaultBody, variableBody} {
		root := t.TempDir()
		script := filepath.Join(root, "lint.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n# shellcheck disable=SC2154\n"+body))
		run(t, ".", "sh", "-n", script)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", script); err != nil {
			t.Fatalf("ShellCheck: %v %s", err, output)
		}
	}
	// Each invocation creates its own root. Resolve the operand before invoking
	// the destructive utility, and require a strict descendant of that root.
	runCase := func(name string, setup func(string) error) (string, int, string) {
		t.Helper()
		root := t.TempDir()
		path := filepath.Join(root, name)
		if name == "-leading" {
			path = name
		}
		absolute := path
		if !filepath.IsAbs(absolute) {
			absolute = filepath.Join(root, path)
		}
		rel, err := filepath.Rel(root, absolute)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
			t.Fatalf("unsafe deletion operand %q relative to %q", absolute, root)
		}
		if setup != nil {
			if err := setup(absolute); err != nil {
				t.Fatal(err)
			}
		}
		script := filepath.Join(root, "delete.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "pathToDirectory="+path), "sh", script)
		return string(output), exitCode(err), absolute
	}
	for _, name := range []string{"space name", "tab\tname", "wild*?name", "-leading", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak", "empty", "nested", "hidden"} {
		name := name
		output, status, path := runCase(name, func(path string) error {
			if err := os.Mkdir(path, 0o700); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(path, ".hidden"), []byte("x"), 0o600); err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(path, "file"), []byte("x"), 0o600); err != nil {
				return err
			}
			if name == "nested" {
				return os.Mkdir(filepath.Join(path, "child"), 0o700)
			}
			return nil
		})
		if status != 0 || output != "" {
			t.Fatalf("delete %q: %d %q", name, status, output)
		}
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("target %q survived: %v", path, err)
		}
	}
	if output, status, _ := runCase("missing", nil); status != 0 || output != "" {
		t.Fatalf("missing: %d %q", status, output)
	}
	if output, status, path := runCase("file", func(path string) error { return os.WriteFile(path, []byte("x"), 0o600) }); status != 0 || output != "" {
		t.Fatalf("file %q: %d %q", path, status, output)
	}
	for _, linkKind := range []string{"outside-dir", "inside-dir", "file", "broken", "operand-dir", "operand-file"} {
		root := t.TempDir()
		target := filepath.Join(root, "target")
		outside := filepath.Join(root, "outside")
		mustMkdirAll(t, outside)
		mustWriteFile(t, filepath.Join(outside, "keep"), []byte("keep"))
		mustMkdirAll(t, target)
		mustMkdirAll(t, filepath.Join(target, "inner"))
		mustWriteFile(t, filepath.Join(target, "regular"), []byte("file"))
		operand := target
		var linkTarget string
		switch linkKind {
		case "outside-dir":
			linkTarget = outside
		case "inside-dir":
			linkTarget = filepath.Join(target, "inner")
		case "file":
			linkTarget = filepath.Join(target, "regular")
		case "broken":
			linkTarget = filepath.Join(root, "absent")
		case "operand-dir":
			operand, linkTarget = filepath.Join(root, "operand"), target
		case "operand-file":
			operand, linkTarget = filepath.Join(root, "operand"), filepath.Join(target, "regular")
		}
		if err := os.Symlink(linkTarget, filepath.Join(target, "link")); err != nil {
			t.Fatal(err)
		}
		if strings.HasPrefix(linkKind, "operand-") {
			if err := os.Symlink(linkTarget, operand); err != nil {
				t.Fatal(err)
			}
		}
		rel, err := filepath.Rel(root, operand)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			t.Fatalf("unsafe symlink operand %q", operand)
		}
		script := filepath.Join(root, "delete.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "pathToDirectory="+operand), "sh", script)
		if err != nil || len(output) != 0 {
			t.Fatalf("symlink %s: %v %q", linkKind, err, output)
		}
		if _, err := os.Lstat(operand); !os.IsNotExist(err) {
			t.Fatalf("operand survived: %s %v", linkKind, err)
		}
		if _, err := os.Stat(filepath.Join(outside, "keep")); err != nil {
			t.Fatalf("outside target lost: %v", err)
		}
		if strings.HasPrefix(linkKind, "operand-") {
			if _, err := os.Stat(target); err != nil {
				t.Fatalf("symlink target lost: %v", err)
			}
		}
	}
}

func TestDirectoryDelete48OperandsAndFailures(t *testing.T) {
	variableBody := expandedGeneratedDirectoryDelete48(t, `${pathToDirectory}`)
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	mustMkdirAll(t, bin)
	capture := filepath.Join(root, "args")
	mustWriteFile(t, filepath.Join(bin, "rm"), []byte("#!/bin/sh\nprintf '%s\\0' \"$@\" > \"$CAPTURE\"\n"))
	if err := os.Chmod(filepath.Join(bin, "rm"), 0o700); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ body, value, want string }{
		{expandedGeneratedDirectoryDelete48(t, "/path/to/directory"), "", "/path/to/directory"},
		{variableBody, "", ""},
		{variableBody, "-leading", "./-leading"},
		{variableBody, "space\t*'\"\\;$()[]\nname", "space\t*'\"\\;$()[]\nname"},
	} {
		script := filepath.Join(root, "capture.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+tc.body))
		env := append(os.Environ(), "PATH="+bin, "CAPTURE="+capture, "pathToDirectory="+tc.value)
		output, err := runCommand(root, env, "/bin/sh", script)
		if err != nil || len(output) != 0 {
			t.Fatalf("capture %q: %v %q", tc.value, err, output)
		}
		args, err := os.ReadFile(capture)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte("-rf\x00" + tc.want + "\x00")
		if !bytes.Equal(args, want) {
			t.Fatalf("operands %q = %q, want %q", tc.value, args, want)
		}
	}
	// The remaining invocations use distinct descendants of a fresh fixture.
	realRoot := t.TempDir()
	script := filepath.Join(realRoot, "delete.sh")
	mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
	runSafe := func(target string, env []string, suffix string) (string, int) {
		t.Helper()
		absolute, err := filepath.Abs(target)
		if err != nil {
			t.Fatal(err)
		}
		rel, err := filepath.Rel(realRoot, absolute)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
			t.Fatalf("unsafe operand: %q", target)
		}
		output, err := runCommand(realRoot, append(env, "pathToDirectory="+target), "/bin/sh", "-c", variableBody+suffix)
		return string(output), exitCode(err)
	}
	missing := filepath.Join(realRoot, "missing")
	if output, status := runSafe(missing, []string{"PATH=/nonexistent"}, ""); status != 127 || !strings.Contains(output, "rm") {
		t.Fatalf("utility unavailable: %d %q", status, output)
	}
	if output, status := runSafe(missing, append(os.Environ(), "PATH=/nonexistent"), "true\n"); status != 0 || !strings.Contains(output, "rm") {
		t.Fatalf("later status masking: %d %q", status, output)
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(realRoot, "blocked")
		mustMkdirAll(t, blocked)
		mustWriteFile(t, filepath.Join(blocked, "file"), []byte("x"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runSafe(blocked, os.Environ(), "")
		if err := os.Chmod(blocked, 0o700); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("permission failure: %d %q", status, output)
		}
	}
}

func expandedGeneratedFileDelete49(t *testing.T, choice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration49MigrationName].Body.(string)
	if !ok || body != currentFileDeleteBody {
		t.Fatalf("generated entry-49 body = %#v", snippets[iteration49MigrationName].Body)
	}
	placeholder := `${1|/path/to/file,${pathToFile}|}`
	if strings.Count(body, placeholder) != 1 {
		t.Fatalf("entry-49 placeholder count in %q", body)
	}
	return strings.Replace(body, placeholder, choice, 1)
}

func TestFileDelete49GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	defaultBody := expandedGeneratedFileDelete49(t, "/path/to/file")
	variableBody := expandedGeneratedFileDelete49(t, `${pathToFile}`)
	for _, body := range []string{defaultBody, variableBody} {
		root := t.TempDir()
		script := filepath.Join(root, "lint.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n# shellcheck disable=SC2154\n"+body))
		run(t, ".", "sh", "-n", script)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", script); err != nil {
			t.Fatalf("ShellCheck: %v %s", err, output)
		}
	}
	for _, name := range []string{"space name", "tab\tname", "wild*?name", "-leading", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak", "empty-file", "read-only-file", "missing", "absolute-file", "relative-file"} {
		root := t.TempDir()
		absolute := filepath.Join(root, name)
		operand := absolute
		if name == "-leading" || name == "relative-file" {
			operand = name
		}
		rel, err := filepath.Rel(root, absolute)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
			t.Fatalf("unsafe deletion target %q", absolute)
		}
		if name != "missing" {
			mode := os.FileMode(0o600)
			if name == "read-only-file" {
				mode = 0
			}
			contents := []byte("payload")
			if name == "empty-file" {
				contents = nil
			}
			mustWriteFile(t, absolute, contents)
			if err := os.Chmod(absolute, mode); err != nil {
				t.Fatal(err)
			}
		}
		script := filepath.Join(root, "delete.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "pathToFile="+operand), "sh", script)
		if err != nil || len(output) != 0 {
			t.Fatalf("delete %q: %v %q", name, err, output)
		}
		if _, err := os.Lstat(absolute); !os.IsNotExist(err) {
			t.Fatalf("target %q survived: %v", name, err)
		}
	}
	for _, kind := range []string{"regular", "directory", "broken"} {
		root := t.TempDir()
		link := filepath.Join(root, "link")
		target := filepath.Join(root, "target")
		if kind == "regular" {
			mustWriteFile(t, target, []byte("keep"))
		} else if kind == "directory" {
			mustMkdirAll(t, target)
			mustWriteFile(t, filepath.Join(target, "keep"), []byte("keep"))
		}
		if err := os.Symlink(target, link); err != nil {
			t.Fatal(err)
		}
		rel, err := filepath.Rel(root, link)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			t.Fatalf("unsafe symlink operand %q", link)
		}
		script := filepath.Join(root, "delete.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "pathToFile="+link), "sh", script)
		if err != nil || len(output) != 0 {
			t.Fatalf("symlink %s: %v %q", kind, err, output)
		}
		if _, err := os.Lstat(link); !os.IsNotExist(err) {
			t.Fatalf("link survived: %v", err)
		}
		if kind != "broken" {
			if _, err := os.Stat(target); err != nil {
				t.Fatalf("symlink target removed: %v", err)
			}
		}
	}
	for _, kind := range []string{"empty-dir", "nonempty-dir"} {
		root := t.TempDir()
		path := filepath.Join(root, kind)
		mustMkdirAll(t, path)
		if kind == "nonempty-dir" {
			mustWriteFile(t, filepath.Join(path, "keep"), []byte("keep"))
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
			t.Fatalf("unsafe directory operand %q", path)
		}
		script := filepath.Join(root, "delete.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "pathToFile="+path), "sh", script)
		if exitCode(err) == 0 || len(output) == 0 {
			t.Fatalf("directory %s: %v %q", kind, err, output)
		}
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("directory lost: %v", err)
		}
	}
	if _, err := exec.LookPath("mkfifo"); err == nil {
		root := t.TempDir()
		path := filepath.Join(root, "fifo")
		if output, err := runCommand(root, nil, "mkfifo", path); err != nil {
			t.Fatalf("mkfifo: %v %q", err, output)
		}
		rel, err := filepath.Rel(root, path)
		if err != nil || rel != "fifo" {
			t.Fatalf("unsafe FIFO operand %q", path)
		}
		script := filepath.Join(root, "delete.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "pathToFile="+path), "sh", script)
		if err != nil || len(output) != 0 {
			t.Fatalf("FIFO: %v %q", err, output)
		}
		if _, err := os.Lstat(path); !os.IsNotExist(err) {
			t.Fatalf("FIFO survived: %v", err)
		}
	}
}

func TestFileDelete49OperandsAndFailures(t *testing.T) {
	variableBody := expandedGeneratedFileDelete49(t, `${pathToFile}`)
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	mustMkdirAll(t, bin)
	capture := filepath.Join(root, "args")
	fakeRM := filepath.Join(bin, "rm")
	mustWriteFile(t, fakeRM, []byte("#!/bin/sh\nprintf '%s\\0' \"$@\" > \"$CAPTURE\"\n"))
	if err := os.Chmod(fakeRM, 0o700); err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ body, value, want string }{
		{expandedGeneratedFileDelete49(t, "/path/to/file"), "", "/path/to/file"},
		{variableBody, "", ""},
		{variableBody, "-leading", "./-leading"},
		{variableBody, "space\t*'\"\\;$()[]\nname", "space\t*'\"\\;$()[]\nname"},
	} {
		script := filepath.Join(root, "capture.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+tc.body))
		env := append(os.Environ(), "PATH="+bin, "CAPTURE="+capture, "pathToFile="+tc.value)
		output, err := runCommand(root, env, "/bin/sh", script)
		if err != nil || len(output) != 0 {
			t.Fatalf("capture %q: %v %q", tc.value, err, output)
		}
		args, err := os.ReadFile(capture)
		if err != nil {
			t.Fatal(err)
		}
		want := []byte("-f\x00" + tc.want + "\x00")
		if !bytes.Equal(args, want) {
			t.Fatalf("operands %q = %q, want %q", tc.value, args, want)
		}
	}
	// The historical variable choice splits a space-containing pathname.
	historical := strings.Replace(v6FileDeleteBody, `${1|/path/to/file,${pathToFile}|}`, `${pathToFile}`, 1)
	script := filepath.Join(root, "historical-capture.sh")
	mustWriteFile(t, script, []byte("#!/bin/sh\n"+historical))
	output, err := runCommand(root, append(os.Environ(), "PATH="+bin, "CAPTURE="+capture, "pathToFile=two words"), "/bin/sh", script)
	if err != nil || len(output) != 0 {
		t.Fatalf("historical capture: %v %q", err, output)
	}
	args, err := os.ReadFile(capture)
	if err != nil {
		t.Fatal(err)
	}
	if want := []byte("-f\x00two\x00words\x00"); !bytes.Equal(args, want) {
		t.Fatalf("historical splitting = %q, want %q", args, want)
	}

	// Every real invocation below validates a strict descendant of its fixture.
	realRoot := t.TempDir()
	realScript := filepath.Join(realRoot, "delete.sh")
	mustWriteFile(t, realScript, []byte("#!/bin/sh\n"+variableBody))
	runSafe := func(target string, env []string, later string) (string, int) {
		t.Helper()
		absolute, err := filepath.Abs(target)
		if err != nil {
			t.Fatal(err)
		}
		rel, err := filepath.Rel(realRoot, absolute)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) || filepath.IsAbs(rel) {
			t.Fatalf("unsafe operand %q", target)
		}
		output, err := runCommand(realRoot, append(env, "pathToFile="+target), "/bin/sh", "-c", variableBody+later)
		return string(output), exitCode(err)
	}
	missing := filepath.Join(realRoot, "missing")
	if output, status := runSafe(missing, []string{"PATH=/nonexistent"}, ""); status != 127 || !strings.Contains(output, "rm") {
		t.Fatalf("missing utility: %d %q", status, output)
	}
	if output, status := runSafe(missing, []string{"PATH=/nonexistent"}, "true\n"); status != 0 || !strings.Contains(output, "rm") {
		t.Fatalf("later status masking: %d %q", status, output)
	}
	if os.Geteuid() != 0 {
		blockedParent := filepath.Join(realRoot, "blocked")
		mustMkdirAll(t, blockedParent)
		blockedFile := filepath.Join(blockedParent, "file")
		mustWriteFile(t, blockedFile, []byte("x"))
		if err := os.Chmod(blockedParent, 0o500); err != nil {
			t.Fatal(err)
		}
		output, status := runSafe(blockedFile, os.Environ(), "")
		if err := os.Chmod(blockedParent, 0o700); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("parent permission failure: %d %q", status, output)
		}
		if _, err := os.Stat(blockedFile); err != nil {
			t.Fatalf("blocked file lost: %v", err)
		}
	}
}

func expandedGeneratedFileRead50(t *testing.T, choice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets[iteration50MigrationName].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentFileReadBody) {
		t.Fatalf("generated entry-50 body = %#v", snippets[iteration50MigrationName].Body)
	}
	lines := make([]string, len(body))
	for i, value := range body {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d = %#v", i, value)
		}
		lines[i] = line
	}
	joined := strings.Join(lines, "\n")
	placeholder := `${1|/path/ro/file,"${inputFile}"|}`
	if strings.Count(joined, placeholder) != 1 {
		t.Fatalf("entry-50 placeholder in %q", joined)
	}
	if choice != "/path/ro/file" && choice != `"${inputFile}"` {
		t.Fatalf("unexpected choice %q", choice)
	}
	joined = strings.Replace(joined, placeholder, choice, 1)
	return strings.ReplaceAll(joined, `\$`, `$`)
}

func TestFileRead50GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	defaultBody := expandedGeneratedFileRead50(t, "/path/ro/file")
	variableBody := expandedGeneratedFileRead50(t, `"${inputFile}"`)
	wantDefault := "file_path=/path/ro/file\nif [ -d \"$file_path\" ]; then\n\tprintf 'cannot read directory: %s\\n' \"$file_path\" >&2\n\tfalse\nelse\n\twhile IFS= read -r line || [ -n \"$line\" ]; do\n\t\tprintf '%s\\n' \"$line\"\n\tdone < \"$file_path\"\nfi\n"
	if defaultBody != wantDefault {
		t.Fatalf("default expansion = %q, want %q", defaultBody, wantDefault)
	}
	if variableBody != strings.Replace(wantDefault, "file_path=/path/ro/file", `file_path="${inputFile}"`, 1) {
		t.Fatalf("variable expansion = %q", variableBody)
	}
	for _, body := range []string{defaultBody, variableBody} {
		root := t.TempDir()
		script := filepath.Join(root, "lint.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n# shellcheck disable=SC2154\n"+body))
		run(t, ".", "sh", "-n", script)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", script); err != nil {
			t.Fatalf("ShellCheck: %v %q", err, output)
		}
	}
	contents := []struct {
		name        string
		input, want []byte
	}{
		{"empty", nil, nil},
		{"one-final-newline", []byte("one\n"), []byte("one\n")},
		{"one-no-newline", []byte("one"), []byte("one\n")},
		{"multiple", []byte("first\nsecond\n"), []byte("first\nsecond\n")},
		{"blank-lines", []byte("\n\nend\n\n"), []byte("\n\nend\n\n")},
		{"spaces-tabs", []byte("  lead\ttrail  \n"), []byte("  lead\ttrail  \n")},
		{"metacharacters", []byte("\\ * ? ' \" $ ; ( ) [ ]\n"), []byte("\\ * ? ' \" $ ; ( ) [ ]\n")},
		{"utf8", []byte("سلام 世界 🌍\n"), []byte("سلام 世界 🌍\n")},
		{"carriage-return", []byte("a\rb\n"), []byte("a\rb\n")},
		{"trailing-newlines", []byte("a\n\n\n"), []byte("a\n\n\n")},
		{"no-final-multiline", []byte("a\nlast"), []byte("a\nlast\n")},
		{"echo-option", []byte("-n\n"), []byte("-n\n")},
	}
	for _, tc := range contents {
		root := t.TempDir()
		path := filepath.Join(root, "input")
		mustWriteFile(t, path, tc.input)
		script := filepath.Join(root, "read.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		output, err := runCommand(root, append(os.Environ(), "inputFile="+path), "sh", script)
		if err != nil || !bytes.Equal(output, tc.want) {
			t.Fatalf("contents %s: %v got %q want %q", tc.name, err, output, tc.want)
		}
	}
	for _, name := range []string{"space name", "tab\tname", "wild*?name", "-leading", "-", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak", "relative-file"} {
		root := t.TempDir()
		path := filepath.Join(root, name)
		mustWriteFile(t, path, []byte("from-file\n"))
		operand := path
		if name == "-leading" || name == "-" || name == "relative-file" {
			operand = name
		}
		script := filepath.Join(root, "read.sh")
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+variableBody))
		command := exec.Command("sh", script)
		command.Dir = root
		command.Env = append(os.Environ(), "inputFile="+operand)
		command.Stdin = strings.NewReader("from-stdin\n")
		output, err := command.CombinedOutput()
		if err != nil || string(output) != "from-file\n" {
			t.Fatalf("pathname %q: %v %q", name, err, output)
		}
	}
}

func TestFileRead50SpecialPathsAndStatus(t *testing.T) {
	body := expandedGeneratedFileRead50(t, `"${inputFile}"`)
	root := t.TempDir()
	runCase := func(path string, later string, environment []string) ([]byte, int) {
		t.Helper()
		output, err := runCommand(root, append(environment, "inputFile="+path), "/bin/sh", "-c", body+later)
		return output, exitCode(err)
	}
	regular := filepath.Join(root, "regular")
	mustWriteFile(t, regular, []byte("linked content\n"))
	link := filepath.Join(root, "link")
	if err := os.Symlink(regular, link); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(link, "", os.Environ()); status != 0 || string(output) != "linked content\n" {
		t.Fatalf("symlink: %d %q", status, output)
	}
	broken := filepath.Join(root, "broken")
	if err := os.Symlink(filepath.Join(root, "absent"), broken); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{broken, filepath.Join(root, "missing"), ""} {
		output, status := runCase(path, "", os.Environ())
		if status == 0 || len(output) == 0 {
			t.Fatalf("missing/empty %q: %d %q", path, status, output)
		}
	}
	directory := filepath.Join(root, "directory")
	mustMkdirAll(t, directory)
	if output, status := runCase(directory, "", os.Environ()); status == 0 || !bytes.Contains(output, []byte("cannot read directory")) {
		t.Fatalf("directory: %d %q", status, output)
	}
	directoryLink := filepath.Join(root, "directory-link")
	if err := os.Symlink(directory, directoryLink); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(directoryLink, "", os.Environ()); status == 0 || !bytes.Contains(output, []byte("cannot read directory")) {
		t.Fatalf("directory symlink: %d %q", status, output)
	}
	if output, status := runCase(filepath.Join(root, "missing"), "true\n", os.Environ()); status != 0 || len(output) == 0 {
		t.Fatalf("later status masking: %d %q", status, output)
	}
	// This snippet uses shell builtins for reading and printing; removing PATH
	// therefore does not make a fixed external read utility unavailable.
	if output, status := runCase(regular, "", []string{"PATH=/nonexistent"}); status != 0 || string(output) != "linked content\n" {
		t.Fatalf("no external utility: %d %q", status, output)
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "blocked")
		mustWriteFile(t, blocked, []byte("secret\n"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runCase(blocked, "", os.Environ())
		if err := os.Chmod(blocked, 0o600); err != nil {
			t.Fatal(err)
		}
		if status == 0 || len(output) == 0 {
			t.Fatalf("unreadable file: %d %q", status, output)
		}
	}
}

func expandedGeneratedFileSearch51(t *testing.T, pathChoice, depth, pattern string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets["filesystem.file-search"].Body.(string)
	if !ok || body != currentFileSearchBody {
		t.Fatalf("generated entry-51 body = %#v", snippets["filesystem.file-search"].Body)
	}
	if pathChoice != "/path/to/search" && pathChoice != `${pathToSearch}` {
		t.Fatalf("unexpected path choice %q", pathChoice)
	}
	if !strings.Contains("|1|2|3|4|5|6|7|8|9|", "|"+depth+"|") {
		t.Fatalf("unexpected depth choice %q", depth)
	}
	for _, substitution := range [][2]string{
		{`${1:result}`, "result"},
		{`${2|/path/to/search,${pathToSearch}|}`, pathChoice},
		{`${3|1,2,3,4,5,6,7,8,9|}`, depth},
		{`${4:criteria}`, pattern},
	} {
		if strings.Count(body, substitution[0]) != 1 {
			t.Fatalf("placeholder %q count in %q", substitution[0], body)
		}
		body = strings.Replace(body, substitution[0], substitution[1], 1)
	}
	return strings.ReplaceAll(body, `\$`, `$`)
}

func TestFileSearch51GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	requireCommand(t, "find")
	requireCommand(t, "grep")
	for _, pathChoice := range []string{"/path/to/search", `${pathToSearch}`} {
		for _, depth := range []string{"1", "2", "3", "4", "5", "6", "7", "8", "9"} {
			body := expandedGeneratedFileSearch51(t, pathChoice, depth, "criteria")
			root := t.TempDir()
			script := filepath.Join(root, "lint.sh")
			mustWriteFile(t, script, []byte("#!/bin/sh\n# shellcheck disable=SC2154 # pathToSearch is supplied by the caller\n"+body+`printf '%s' "$result"`+"\n"))
			run(t, ".", "sh", "-n", script)
			if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", script); err != nil {
				t.Fatalf("ShellCheck: %v %q", err, output)
			}
		}
	}
	root := t.TempDir()
	runSearch := func(path, pattern string, depth string, env []string) (string, int) {
		t.Helper()
		body := expandedGeneratedFileSearch51(t, `${pathToSearch}`, depth, `$searchPattern`)
		command := exec.Command("/bin/sh", "-c", body+"search_status=$?\nprintf '%s' \"$result\"\nexit \"$search_status\"\n")
		command.Dir = root
		command.Env = append(append(os.Environ(), "pathToSearch="+path, "searchPattern="+pattern), env...)
		output, err := command.CombinedOutput()
		return string(output), exitCode(err)
	}
	for _, name := range []string{"space name", "tab\tname", "wild*?name", "-leading", "-", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak"} {
		path := filepath.Join(root, name)
		mustWriteFile(t, path, []byte("needle\n"))
		operand := path
		if name == "-leading" || name == "-" {
			operand = name
		}
		output, status := runSearch(operand, "needle", "1", nil)
		if status != 0 || output != "needle" {
			t.Fatalf("pathname %q: status %d, output %q", name, status, output)
		}
	}
	regular := filepath.Join(root, "pattern-input")
	mustWriteFile(t, regular, []byte("a b\na\tb\n-dot\nstar*\n[abc]\nanchor\nback\\slash\nsingle'quote\ndouble\"quote\ndollar$sign\nsemi;colon\nround(parens)\n\n"))
	findBatchProbe := exec.Command("find", regular, "-maxdepth", "0", "-type", "f", "-exec", "sh", "-c", "exit 17", "sh", "{}", "+")
	if output, err := findBatchProbe.CombinedOutput(); exitCode(err) != 1 || len(output) != 0 {
		t.Fatalf("find batch-shell failure propagation: %v %q", err, output)
	}
	patterns := []string{"a b", "a\tb", "-dot", "", ".", "star*", "[abc]", "^anchor$", `back\\slash`, "single'quote", `double"quote`, "dollar$sign", "semi;colon", `round(parens)`, "a\nb"}
	for _, pattern := range patterns {
		output, status := runSearch(regular, pattern, "1", nil)
		if status != 0 || output == "" && pattern != "" {
			t.Fatalf("pattern %q: status %d, output %q", pattern, status, output)
		}
	}
	if output, status := runSearch(regular, "absent-pattern", "1", nil); status != 1 || output != "" {
		t.Fatalf("no match: %d %q", status, output)
	}
	if output, status := runSearch(regular, "[", "1", nil); status != 2 || !strings.Contains(output, "grep:") {
		t.Fatalf("malformed BRE: %d %q", status, output)
	}
	if output, status := runSearch("missing-file", "needle", "1", nil); status != 2 || !strings.Contains(output, "find:") {
		t.Fatalf("nonexistent: %d %q", status, output)
	}
	if output, status := runSearch(".", "needle", "1", nil); status != 0 || !strings.Contains(output, "needle") {
		t.Fatalf("directory search: %d %q", status, output)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(regular, link); err != nil {
		t.Fatal(err)
	}
	if output, status := runSearch(link, "anchor", "1", nil); status != 1 || output != "" {
		t.Fatalf("symlink operand: %d %q", status, output)
	}
	broken := filepath.Join(root, "broken")
	if err := os.Symlink(filepath.Join(root, "absent-target"), broken); err != nil {
		t.Fatal(err)
	}
	if output, status := runSearch(broken, "anchor", "1", nil); status != 1 || output != "" {
		t.Fatalf("broken symlink: %d %q", status, output)
	}
	empty := filepath.Join(root, "empty")
	mustWriteFile(t, empty, nil)
	if output, status := runSearch(empty, "needle", "1", nil); status != 1 || output != "" {
		t.Fatalf("empty file: %d %q", status, output)
	}
	emptyDirectory := filepath.Join(root, "empty-directory")
	mustMkdirAll(t, emptyDirectory)
	if output, status := runSearch(emptyDirectory, "needle", "1", nil); status != 1 || output != "" {
		t.Fatalf("zero selected regular files: %d %q", status, output)
	}
	tree := filepath.Join(root, "tree")
	mustMkdirAll(t, filepath.Join(tree, "child"))
	mustWriteFile(t, filepath.Join(tree, "top"), []byte("top-marker\n"))
	mustWriteFile(t, filepath.Join(tree, "child", "nested"), []byte("nested-marker\n"))
	if output, status := runSearch(tree, "nested-marker", "1", nil); status != 1 || output != "" {
		t.Fatalf("depth one: %d %q", status, output)
	}
	if output, status := runSearch(tree, "nested-marker", "2", nil); status != 0 || !strings.HasSuffix(output, ":nested-marker") {
		t.Fatalf("depth two: %d %q", status, output)
	}
	if output, status := runSearch("pattern-input", "anchor", "1", []string{"PATH=/nonexistent"}); status != 2 || !strings.Contains(output, "mktemp:") {
		t.Fatalf("utility unavailable: %d %q", status, output)
	}
	limitedPath := filepath.Join(root, "limited-bin")
	mustMkdirAll(t, limitedPath)
	for _, utility := range []string{"mktemp", "find", "sh", "rm"} {
		found, err := exec.LookPath(utility)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(found, filepath.Join(limitedPath, utility)); err != nil {
			t.Fatal(err)
		}
	}
	if output, status := runSearch("pattern-input", "anchor", "1", []string{"PATH=" + limitedPath}); status != 2 || !strings.Contains(output, "grep:") {
		t.Fatalf("grep unavailable: %d %q", status, output)
	}
	interceptPath := filepath.Join(root, "intercept-bin")
	mustMkdirAll(t, interceptPath)
	for _, utility := range []string{"find", "sh", "grep", "rm"} {
		found, err := exec.LookPath(utility)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(found, filepath.Join(interceptPath, utility)); err != nil {
			t.Fatal(err)
		}
	}
	fakeMktemp := filepath.Join(interceptPath, "mktemp")
	mustWriteFile(t, fakeMktemp, []byte("#!/bin/sh\nprintf '%s\\n' \"$SHELLMAN_STATUS_TARGET\"\n"))
	if err := os.Chmod(fakeMktemp, 0o755); err != nil {
		t.Fatal(err)
	}
	fakeCatForAppend := filepath.Join(interceptPath, "cat")
	mustWriteFile(t, fakeCatForAppend, []byte("#!/bin/sh\n: > \"$SHELLMAN_CAT_MARKER\"\n"))
	if err := os.Chmod(fakeCatForAppend, 0o755); err != nil {
		t.Fatal(err)
	}
	statusSink := filepath.Join(root, "status-is-directory")
	mustMkdirAll(t, statusSink)
	catMarker := filepath.Join(root, "cat-was-called")
	if output, status := runSearch(regular, "anchor", "1", []string{"PATH=" + interceptPath, "SHELLMAN_STATUS_TARGET=" + statusSink, "SHELLMAN_CAT_MARKER=" + catMarker}); status != 2 || output == "" {
		t.Fatalf("status append failure: %d %q", status, output)
	}
	if _, err := os.Stat(catMarker); !os.IsNotExist(err) {
		t.Fatalf("find did not stop status-file aggregation: cat marker stat = %v", err)
	}
	mustWriteFile(t, fakeMktemp, []byte("#!/bin/sh\nexit 7\n"))
	if output, status := runSearch(regular, "anchor", "1", []string{"PATH=" + interceptPath}); status != 2 || output != "" {
		t.Fatalf("temporary-file creation failure: %d %q", status, output)
	}
	readInterceptPath := filepath.Join(root, "read-intercept-bin")
	mustMkdirAll(t, readInterceptPath)
	for _, utility := range []string{"mktemp", "find", "sh", "grep", "rm"} {
		found, err := exec.LookPath(utility)
		if err != nil {
			t.Fatal(err)
		}
		if err := os.Symlink(found, filepath.Join(readInterceptPath, utility)); err != nil {
			t.Fatal(err)
		}
	}
	fakeCat := filepath.Join(readInterceptPath, "cat")
	mustWriteFile(t, fakeCat, []byte("#!/bin/sh\nexit 7\n"))
	if err := os.Chmod(fakeCat, 0o755); err != nil {
		t.Fatal(err)
	}
	if output, status := runSearch(emptyDirectory, "needle", "1", []string{"PATH=" + readInterceptPath}); status != 2 || output != "" {
		t.Fatalf("status-file read failure: %d %q", status, output)
	}
	binary := filepath.Join(root, "binary")
	mustWriteFile(t, binary, []byte{'a', 0, 'a', 'n', 'c', 'h', 'o', 'r', '\n'})
	if output, status := runSearch(binary, "anchor", "1", nil); status != 0 || output == "" {
		t.Fatalf("GNU grep binary fixture: %d %q", status, output)
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "blocked")
		mustWriteFile(t, blocked, []byte("secret\n"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runSearch(blocked, "secret", "1", nil)
		if err := os.Chmod(blocked, 0o600); err != nil {
			t.Fatal(err)
		}
		if status == 0 || !strings.Contains(output, "grep:") {
			t.Fatalf("unreadable file: %d %q", status, output)
		}
	}
}

func expandedGeneratedFileWriteMultilineSudo52(t *testing.T, pathChoice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets["filesystem.file-write-multiline-sudo"].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentFileWriteMultilineSudoBody) {
		t.Fatalf("generated entry-52 body = %#v", snippets["filesystem.file-write-multiline-sudo"].Body)
	}
	if pathChoice != "/path/to/file" && pathChoice != "${filePath}" {
		t.Fatalf("unexpected path choice %q", pathChoice)
	}
	lines := make([]string, len(body))
	for index, value := range body {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d = %#v", index, value)
		}
		lines[index] = line
	}
	joined := strings.Join(lines, "\n")
	placeholder := `${1|/path/to/file,${filePath}|}`
	if strings.Count(joined, placeholder) != 1 {
		t.Fatalf("destination placeholder count in %q", joined)
	}
	joined = strings.Replace(joined, placeholder, pathChoice, 1)
	return strings.ReplaceAll(joined, `\$`, `$`)
}

func TestFileWriteMultilineSudo52GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	requireCommand(t, "script")
	requireCommand(t, "tee")
	for _, choice := range []string{"/path/to/file", "${filePath}"} {
		body := expandedGeneratedFileWriteMultilineSudo52(t, choice)
		root := t.TempDir()
		scriptPath := filepath.Join(root, "lint.sh")
		mustWriteFile(t, scriptPath, []byte("#!/bin/sh\n# shellcheck disable=SC2154 # filePath is supplied by the caller\n"+body))
		run(t, ".", "sh", "-n", scriptPath)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", scriptPath); err != nil {
			t.Fatalf("ShellCheck: %v %q", err, output)
		}
	}

	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	mustMkdirAll(t, bin)
	mustExecutable(t, filepath.Join(bin, "sudo"), "#!/bin/sh\n[ \"$#\" -eq 2 ] && [ \"$1\" = tee ] || exit 98\nprintf '%s' \"$2\" > \"$SHELLMAN_SUDO_PATH_LOG\"\n[ \"${SHELLMAN_SUDO_FAIL:-}\" = 1 ] && exit 37\nif [ \"${SHELLMAN_SUDO_PASSWORD:-}\" = 1 ]; then IFS= read -r password </dev/tty || exit 38; [ \"$password\" = secret ] || exit 39; fi\nexec \"$@\"\n")
	teePath, err := exec.LookPath("tee")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(teePath, filepath.Join(bin, "tee")); err != nil {
		t.Fatal(err)
	}
	logPath := filepath.Join(root, "sudo-path")
	runCase := func(path, input string, extra []string, later string) (string, int) {
		t.Helper()
		body := expandedGeneratedFileWriteMultilineSudo52(t, "${filePath}")
		scriptPath := filepath.Join(root, "write.sh")
		mustWriteFile(t, scriptPath, []byte("#!/bin/sh\n"+body+later))
		quotedScript := "'" + strings.ReplaceAll(scriptPath, "'", "'\\''") + "'"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, "script", "-q", "-e", "-c", "/bin/sh "+quotedScript, "/dev/null")
		command.Dir = root
		command.Stdin = strings.NewReader(input)
		command.Env = append(append(os.Environ(), "PATH="+bin, "filePath="+path, "SHELLMAN_SUDO_PATH_LOG="+logPath), extra...)
		output, err := command.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("PTY command timed out: %v %q", ctx.Err(), output)
		}
		return string(output), exitCode(err)
	}
	end := "\x04"
	for _, tc := range []struct{ name, input, want string }{
		{"empty", end, ""},
		{"one empty line", "\n" + end, "\n"},
		{"multiple blank lines", "\n\n" + end, "\n\n"},
		{"final newline", "alpha\nbeta\n" + end, "alpha\nbeta\n"},
		{"no final newline", "alpha\nbeta" + end + end, "alpha\nbeta"},
		{"literal delimiter lines", "EOF\n_SHELLMAN_STATUS\nEND\n" + end, "EOF\n_SHELLMAN_STATUS\nEND\n"},
		{"CR maps to LF", "a\rb\n" + end, "a\nb\n"},
		{"CRLF means two terminal line breaks", "a\r\n" + end, "a\n\n"},
	} {
		path := filepath.Join(root, "content-output")
		output, status := runCase(path, tc.input, nil, "")
		got, err := os.ReadFile(path)
		if err != nil || status != 0 || string(got) != tc.want {
			t.Fatalf("%s: status %d output %q bytes %q error %v", tc.name, status, output, got, err)
		}
	}
	marker := filepath.Join(root, "MUST-NOT-EXIST")
	malicious := "$(touch " + marker + ")\n`touch " + marker + "`\n; touch " + marker + "\nEOF\n_SHELLMAN_STATUS\nquotes ' \" and back\\slash ; | & * ? [ ] ( ) $\n"
	path := filepath.Join(root, "malicious-output")
	if output, status := runCase(path, malicious+end, nil, ""); status != 0 {
		t.Fatalf("malicious-looking data: %d %q", status, output)
	}
	if got, err := os.ReadFile(path); err != nil || string(got) != malicious {
		t.Fatalf("malicious-looking bytes: %q %v", got, err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("content executed: %v", err)
	}

	for _, name := range []string{"absolute", "relative", "space name", "tab\tname", "wild*?name", "-leading", "-", "single'quote", `double"quote`, `back\slash`, "semi;colon", "dollar$sign", "round(parens)", "square[brackets]", "line\nbreak"} {
		operand := filepath.Join(root, name)
		if name == "relative" || name == "-leading" || name == "-" {
			operand = name
		}
		output, status := runCase(operand, "path payload\n"+end, nil, "")
		if status != 0 {
			t.Fatalf("pathname %q: %d %q", name, status, output)
		}
		wantArgument := operand
		if strings.HasPrefix(operand, "-") {
			wantArgument = "./" + operand
		}
		argument, err := os.ReadFile(logPath)
		if err != nil || string(argument) != wantArgument {
			t.Fatalf("sudo tee argument %q = %q, %v", name, argument, err)
		}
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || string(got) != "path payload\n" {
			t.Fatalf("pathname write %q: %q %v", name, got, err)
		}
	}
	if output, status := runCase("", "ignored\n"+end, nil, ""); status != 2 || !strings.Contains(output, "destination is empty") {
		t.Fatalf("empty destination: %d %q", status, output)
	}
	if output, status := runCase(filepath.Join(root, "password-output"), "secret\ncontent\n"+end, []string{"SHELLMAN_SUDO_PASSWORD=1"}, ""); status != 0 {
		t.Fatalf("sudo terminal authentication: %d %q", status, output)
	}
	if got, err := os.ReadFile(filepath.Join(root, "password-output")); err != nil || string(got) != "content\n" {
		t.Fatalf("password consumed as content: %q %v", got, err)
	}
	existing := filepath.Join(root, "existing")
	mustWriteFile(t, existing, []byte("old data\n"))
	if output, status := runCase(existing, "new\n"+end, nil, ""); status != 0 {
		t.Fatalf("overwrite existing: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "new\n" {
		t.Fatalf("overwrite bytes: %q %v", got, err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(existing, link); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(link, "through-link\n"+end, nil, ""); status != 0 {
		t.Fatalf("symlink: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "through-link\n" {
		t.Fatalf("symlink target: %q %v", got, err)
	}
	for _, badPath := range []string{filepath.Join(root, "missing-parent", "file"), root} {
		if output, status := runCase(badPath, "cannot-write\n"+end, nil, ""); status == 0 || output == "" {
			t.Fatalf("invalid destination %q: %d %q", badPath, status, output)
		}
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "unwritable")
		mustWriteFile(t, blocked, []byte("preserve me\n"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runCase(blocked, "cannot-write\n"+end, nil, "")
		if err := os.Chmod(blocked, 0o600); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("permission failure: %d %q", status, output)
		}
	}
	if output, status := runCase(existing, "not-written\n"+end, []string{"SHELLMAN_SUDO_FAIL=1"}, "true\n"); status != 37 {
		t.Fatalf("sudo failure or later status masking: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "through-link\n" {
		t.Fatalf("sudo failure changed destination: %q %v", got, err)
	}
	if output, status := runCase(existing, "cancel\n\x03", nil, ""); status == 0 {
		t.Fatalf("interrupted input: %d %q", status, output)
	}
	if err := os.Remove(filepath.Join(bin, "tee")); err != nil {
		t.Fatal(err)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\nprintf 'partial' > \"$1\"\nexit 29\n")
	if output, status := runCase(existing, "full content\n"+end, nil, "true\n"); status != 29 {
		t.Fatalf("writer failure or later status masking: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "partial" {
		t.Fatalf("partial write: %q %v", got, err)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\nkill -TERM \"$$\"\n")
	if output, status := runCase(existing, "interrupted\n"+end, nil, ""); status == 0 {
		t.Fatalf("interrupted writer: %d %q", status, output)
	}
	if err := os.Remove(filepath.Join(bin, "tee")); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(existing, "missing tee\n"+end, nil, ""); status == 0 {
		t.Fatalf("missing writer: %d %q", status, output)
	}
	if err := os.Remove(filepath.Join(bin, "sudo")); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(existing, "missing sudo\n"+end, nil, ""); status == 0 {
		t.Fatalf("missing sudo: %d %q", status, output)
	}
	// A non-terminal stream must fail before sudo is invoked.
	nonTerminal := exec.Command("/bin/sh", filepath.Join(root, "write.sh"))
	nonTerminal.Env = append(os.Environ(), "filePath="+existing, "PATH="+bin, "SHELLMAN_SUDO_PATH_LOG="+logPath)
	nonTerminal.Stdin = strings.NewReader("not interactive\n")
	if output, err := nonTerminal.CombinedOutput(); exitCode(err) != 2 || !strings.Contains(string(output), "interactive terminal input is required") {
		t.Fatalf("non-terminal input: %d %q", exitCode(err), output)
	}
	quotedScript := "'" + strings.ReplaceAll(filepath.Join(root, "write.sh"), "'", "'\\''") + "'"
	closedInput := exec.Command("/bin/sh", "-c", "exec 0<&-; exec /bin/sh "+quotedScript)
	closedInput.Env = append(os.Environ(), "filePath="+existing, "PATH="+bin, "SHELLMAN_SUDO_PATH_LOG="+logPath)
	if output, err := closedInput.CombinedOutput(); exitCode(err) != 2 || !strings.Contains(string(output), "interactive terminal input is required") {
		t.Fatalf("closed input: %d %q", exitCode(err), output)
	}
}

func expandedGeneratedFileWriteMultiline53(t *testing.T, pathChoice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets["filesystem.file-write-multiline"].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentFileWriteMultilineBody) {
		t.Fatalf("generated entry-53 body = %#v", snippets["filesystem.file-write-multiline"].Body)
	}
	if pathChoice != "/path/to/file" && pathChoice != "${filePath}" {
		t.Fatalf("unexpected path choice %q", pathChoice)
	}
	lines := make([]string, len(body))
	for index, value := range body {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d = %#v", index, value)
		}
		lines[index] = line
	}
	joined := strings.Join(lines, "\n")
	placeholder := `${1|/path/to/file,${filePath}|}`
	if strings.Count(joined, placeholder) != 1 {
		t.Fatalf("destination placeholder count in %q", joined)
	}
	return strings.ReplaceAll(strings.Replace(joined, placeholder, pathChoice, 1), `\$`, `$`)
}

func TestFileWriteMultiline53GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	requireCommand(t, "script")
	requireCommand(t, "tee")
	for _, choice := range []string{"/path/to/file", "${filePath}"} {
		body := expandedGeneratedFileWriteMultiline53(t, choice)
		lint := filepath.Join(t.TempDir(), "lint.sh")
		mustWriteFile(t, lint, []byte("#!/bin/sh\n# shellcheck disable=SC2154 # filePath is supplied by the caller\n"+body))
		run(t, ".", "sh", "-n", lint)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", lint); err != nil {
			t.Fatalf("ShellCheck: %v %q", err, output)
		}
	}
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	mustMkdirAll(t, bin)
	teePath, err := exec.LookPath("tee")
	if err != nil {
		t.Fatal(err)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\n[ \"$#\" -eq 1 ] || exit 98\nprintf '%s' \"$1\" > \"$SHELLMAN_TEE_PATH_LOG\"\n[ \"${SHELLMAN_TEE_FAIL:-}\" = 1 ] && exit 37\nexec \"$SHELLMAN_REAL_TEE\" \"$@\"\n")
	logPath := filepath.Join(root, "tee-path")
	runCase := func(path, input string, extra []string, later string) (string, int) {
		t.Helper()
		scriptPath := filepath.Join(root, "write.sh")
		mustWriteFile(t, scriptPath, []byte("#!/bin/sh\n"+expandedGeneratedFileWriteMultiline53(t, "${filePath}")+later))
		quotedScript := "'" + strings.ReplaceAll(scriptPath, "'", "'\\''") + "'"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, "script", "-q", "-e", "-c", "/bin/sh "+quotedScript, "/dev/null")
		command.Dir = root
		command.Stdin = strings.NewReader(input)
		command.Env = append(append(os.Environ(), "PATH="+bin, "filePath="+path, "SHELLMAN_TEE_PATH_LOG="+logPath, "SHELLMAN_REAL_TEE="+teePath), extra...)
		output, err := command.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("PTY command timed out: %v %q", ctx.Err(), output)
		}
		return string(output), exitCode(err)
	}
	end := "\x04"
	for _, tc := range []struct{ name, input, want string }{
		{"empty", end, ""},
		{"one blank line", "\n" + end, "\n"},
		{"multiple blank lines", "\n\n" + end, "\n\n"},
		{"final LF", "alpha\nbeta\n" + end, "alpha\nbeta\n"},
		{"no final LF", "alpha\nbeta" + end + end, "alpha\nbeta"},
		{"delimiter lines", "EOF\n_SHELLMAN_STATUS\nEND\n" + end, "EOF\n_SHELLMAN_STATUS\nEND\n"},
		{"terminal CR", "a\rb\n" + end, "a\nb\n"},
		{"terminal CRLF", "a\r\n" + end, "a\n\n"},
	} {
		path := filepath.Join(root, "content-output")
		output, status := runCase(path, tc.input, nil, "")
		got, err := os.ReadFile(path)
		if status != 0 || err != nil || string(got) != tc.want {
			t.Fatalf("%s: status %d output %q bytes %q error %v", tc.name, status, output, got, err)
		}
	}
	marker := filepath.Join(root, "MUST-NOT-EXIST")
	malicious := "$(touch " + marker + ")\n`touch " + marker + "`\n; touch " + marker + "\nEOF\n_SHELLMAN_STATUS\nquotes ' \" back\\slash $ * ? [ ] ( ) | &\n"
	maliciousPath := filepath.Join(root, "malicious-output")
	if output, status := runCase(maliciousPath, malicious+end, nil, ""); status != 0 {
		t.Fatalf("malicious-looking input: %d %q", status, output)
	}
	if got, err := os.ReadFile(maliciousPath); err != nil || string(got) != malicious {
		t.Fatalf("malicious-looking bytes: %q %v", got, err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("entered input executed: %v", err)
	}
	for _, name := range []string{"absolute", "relative", "space name", "tab\tname", "wild*?name", "-leading", "-", "single'quote", `double"quote`, `back\slash`, "dollar$sign", "semi;colon", "round(parens)", "square[brackets]", "line\nbreak"} {
		operand := filepath.Join(root, name)
		if name == "relative" || name == "-leading" || name == "-" {
			operand = name
		}
		output, status := runCase(operand, "path payload\n"+end, nil, "")
		if status != 0 {
			t.Fatalf("pathname %q: %d %q", name, status, output)
		}
		wantArgument := operand
		if strings.HasPrefix(operand, "-") {
			wantArgument = "./" + operand
		}
		argument, err := os.ReadFile(logPath)
		if err != nil || string(argument) != wantArgument {
			t.Fatalf("tee argument %q = %q, %v", name, argument, err)
		}
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || string(got) != "path payload\n" {
			t.Fatalf("pathname write %q: %q %v", name, got, err)
		}
	}
	if output, status := runCase("", "ignored\n"+end, nil, ""); status != 2 || !strings.Contains(output, "destination is empty") {
		t.Fatalf("empty destination: %d %q", status, output)
	}
	existing := filepath.Join(root, "existing")
	mustWriteFile(t, existing, []byte("old data\n"))
	if output, status := runCase(existing, "new\n"+end, nil, ""); status != 0 {
		t.Fatalf("overwrite: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "new\n" {
		t.Fatalf("overwrite bytes: %q %v", got, err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(existing, link); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(link, "through-link\n"+end, nil, ""); status != 0 {
		t.Fatalf("symlink: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "through-link\n" {
		t.Fatalf("symlink target: %q %v", got, err)
	}
	for _, badPath := range []string{filepath.Join(root, "missing-parent", "file"), root} {
		if output, status := runCase(badPath, "cannot-write\n"+end, nil, ""); status == 0 || output == "" {
			t.Fatalf("invalid destination %q: %d %q", badPath, status, output)
		}
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "unwritable")
		mustWriteFile(t, blocked, []byte("old"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runCase(blocked, "cannot-write\n"+end, nil, "")
		if err := os.Chmod(blocked, 0o600); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("permission failure: %d %q", status, output)
		}
	}
	if output, status := runCase(existing, "not-written\n"+end, []string{"SHELLMAN_TEE_FAIL=1"}, ""); status != 37 {
		t.Fatalf("writer failure status: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "through-link\n" {
		t.Fatalf("failed writer changed destination: %q %v", got, err)
	}
	if output, status := runCase(existing, "cancel\n\x03", nil, ""); status == 0 {
		t.Fatalf("interrupted input: %d %q", status, output)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\nprintf 'partial' > \"$1\"\nexit 29\n")
	if output, status := runCase(existing, "full content\n"+end, nil, ""); status != 29 {
		t.Fatalf("partial writer failure: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "partial" {
		t.Fatalf("partial writer bytes: %q %v", got, err)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\nkill -TERM \"$$\"\n")
	if output, status := runCase(existing, "interrupted writer\n"+end, nil, ""); status == 0 {
		t.Fatalf("interrupted writer: %d %q", status, output)
	}
	if err := os.Remove(filepath.Join(bin, "tee")); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(existing, "missing tee\n"+end, nil, ""); status == 0 {
		t.Fatalf("missing tee: %d %q", status, output)
	}
	if output, status := runCase(existing, "missing tee\n"+end, nil, "printf 'AFTER\\n'\n"); status != 0 || !strings.Contains(output, "AFTER") {
		t.Fatalf("later-command masking or shell continuation: %d %q", status, output)
	}
	// Nonterminal and closed stdin fail before a destination is opened.
	scriptPath := filepath.Join(root, "write.sh")
	mustWriteFile(t, scriptPath, []byte("#!/bin/sh\n"+expandedGeneratedFileWriteMultiline53(t, "${filePath}")))
	for _, closeFD := range []bool{false, true} {
		var command *exec.Cmd
		if closeFD {
			quoted := "'" + strings.ReplaceAll(scriptPath, "'", "'\\''") + "'"
			command = exec.Command("/bin/sh", "-c", "exec 0<&-; exec /bin/sh "+quoted)
		} else {
			command = exec.Command("/bin/sh", scriptPath)
			command.Stdin = strings.NewReader("redirected\n")
		}
		unopened := filepath.Join(root, "must-not-open")
		command.Env = append(os.Environ(), "filePath="+unopened, "PATH="+bin)
		output, err := command.CombinedOutput()
		if exitCode(err) != 2 || !strings.Contains(string(output), "interactive terminal input is required") {
			t.Fatalf("nonterminal or closed input: %d %q", exitCode(err), output)
		}
		if _, err := os.Stat(unopened); !os.IsNotExist(err) {
			t.Fatalf("nonterminal input opened destination: %v", err)
		}
	}
}

func expandedGeneratedFileWrite54(t *testing.T, pathChoice string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets["filesystem.file-write"].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentFileWriteBody) {
		t.Fatalf("generated entry-54 body = %#v", snippets["filesystem.file-write"].Body)
	}
	if pathChoice != "/path/to/file" && pathChoice != "${pathToFile}" {
		t.Fatalf("unexpected path choice %q", pathChoice)
	}
	lines := make([]string, len(body))
	for index, value := range body {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d = %#v", index, value)
		}
		lines[index] = line
	}
	joined := strings.Join(lines, "\n")
	placeholder := `${2|/path/to/file,${pathToFile}|}`
	if strings.Count(joined, placeholder) != 1 {
		t.Fatalf("destination placeholder count in %q", joined)
	}
	return strings.ReplaceAll(strings.Replace(joined, placeholder, pathChoice, 1), `\$`, `$`)
}

func TestFileWrite54GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "shellcheck")
	requireCommand(t, "script")
	requireCommand(t, "tee")
	for _, choice := range []string{"/path/to/file", "${pathToFile}"} {
		body := expandedGeneratedFileWrite54(t, choice)
		lint := filepath.Join(t.TempDir(), "lint.sh")
		mustWriteFile(t, lint, []byte("#!/bin/sh\n# shellcheck disable=SC2154 # pathToFile is supplied by the caller\n"+body))
		run(t, ".", "sh", "-n", lint)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", lint); err != nil {
			t.Fatalf("ShellCheck: %v %q", err, output)
		}
	}
	root := t.TempDir()
	bin := filepath.Join(root, "bin")
	mustMkdirAll(t, bin)
	teePath, err := exec.LookPath("tee")
	if err != nil {
		t.Fatal(err)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\n[ \"$#\" -eq 1 ] || exit 98\nprintf '%s' \"$1\" > \"$SHELLMAN_TEE_PATH_LOG\"\n[ \"${SHELLMAN_TEE_FAIL:-}\" = 1 ] && exit 37\nexec \"$SHELLMAN_REAL_TEE\" \"$@\"\n")
	logPath := filepath.Join(root, "tee-path")
	runCase := func(path, input string, extra []string, later string) (string, int) {
		t.Helper()
		scriptPath := filepath.Join(root, "write.sh")
		mustWriteFile(t, scriptPath, []byte("#!/bin/sh\n"+expandedGeneratedFileWrite54(t, "${pathToFile}")+later))
		quotedScript := "'" + strings.ReplaceAll(scriptPath, "'", "'\\''") + "'"
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		command := exec.CommandContext(ctx, "script", "-q", "-e", "-c", "/bin/sh "+quotedScript, "/dev/null")
		command.Dir = root
		command.Stdin = strings.NewReader(input)
		command.Env = append(append(os.Environ(), "PATH="+bin, "pathToFile="+path, "SHELLMAN_TEE_PATH_LOG="+logPath, "SHELLMAN_REAL_TEE="+teePath), extra...)
		output, err := command.CombinedOutput()
		if ctx.Err() != nil {
			t.Fatalf("PTY command timed out: %v %q", ctx.Err(), output)
		}
		return string(output), exitCode(err)
	}
	end := "\x04"
	for _, tc := range []struct{ name, input, want string }{
		{"empty", end, ""},
		{"one blank line", "\n" + end, "\n"},
		{"multiple blank lines", "\n\n" + end, "\n\n"},
		{"final LF", "alpha\nbeta\n" + end, "alpha\nbeta\n"},
		{"no final LF", "alpha\nbeta" + end + end, "alpha\nbeta"},
		{"delimiter lines", "EOF\n_SHELLMAN_STATUS\nEND\n" + end, "EOF\n_SHELLMAN_STATUS\nEND\n"},
		{"terminal CR", "a\rb\n" + end, "a\nb\n"},
		{"terminal CRLF", "a\r\n" + end, "a\n\n"},
	} {
		path := filepath.Join(root, "content-output")
		output, status := runCase(path, tc.input, nil, "")
		got, err := os.ReadFile(path)
		if status != 0 || err != nil || string(got) != tc.want {
			t.Fatalf("%s: status %d output %q bytes %q error %v", tc.name, status, output, got, err)
		}
	}
	marker := filepath.Join(root, "MUST-NOT-EXIST")
	malicious := "$(touch " + marker + ")\n`touch " + marker + "`\n; touch " + marker + "\nEOF\n_SHELLMAN_STATUS\nquotes ' \" back\\slash $ * ? [ ] ( ) | &\n"
	maliciousPath := filepath.Join(root, "malicious-output")
	if output, status := runCase(maliciousPath, malicious+end, nil, ""); status != 0 {
		t.Fatalf("malicious-looking input: %d %q", status, output)
	}
	if got, err := os.ReadFile(maliciousPath); err != nil || string(got) != malicious {
		t.Fatalf("malicious-looking bytes: %q %v", got, err)
	}
	if _, err := os.Stat(marker); !os.IsNotExist(err) {
		t.Fatalf("entered input executed: %v", err)
	}
	for _, name := range []string{"absolute", "relative", "space name", "tab\tname", "wild*?name", "-leading", "-", "single'quote", `double"quote`, `back\slash`, "dollar$sign", "semi;colon", "round(parens)", "square[brackets]", "line\nbreak"} {
		operand := filepath.Join(root, name)
		if name == "relative" || name == "-leading" || name == "-" {
			operand = name
		}
		output, status := runCase(operand, "path payload\n"+end, nil, "")
		if status != 0 {
			t.Fatalf("pathname %q: %d %q", name, status, output)
		}
		wantArgument := operand
		if strings.HasPrefix(operand, "-") {
			wantArgument = "./" + operand
		}
		argument, err := os.ReadFile(logPath)
		if err != nil || string(argument) != wantArgument {
			t.Fatalf("tee argument %q = %q, %v", name, argument, err)
		}
		got, err := os.ReadFile(filepath.Join(root, name))
		if err != nil || string(got) != "path payload\n" {
			t.Fatalf("pathname write %q: %q %v", name, got, err)
		}
	}
	if err := os.Remove(logPath); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase("", "ignored\n"+end, nil, ""); status != 2 || !strings.Contains(output, "destination is empty") {
		t.Fatalf("empty destination: %d %q", status, output)
	}
	if _, err := os.Stat(logPath); !os.IsNotExist(err) {
		t.Fatalf("empty destination invoked tee: %v", err)
	}
	existing := filepath.Join(root, "existing")
	mustWriteFile(t, existing, []byte("old data\n"))
	if output, status := runCase(existing, "new\n"+end, nil, ""); status != 0 {
		t.Fatalf("overwrite: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "new\n" {
		t.Fatalf("overwrite bytes: %q %v", got, err)
	}
	link := filepath.Join(root, "link")
	if err := os.Symlink(existing, link); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(link, "through-link\n"+end, nil, ""); status != 0 {
		t.Fatalf("symlink: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "through-link\n" {
		t.Fatalf("symlink target: %q %v", got, err)
	}
	for _, badPath := range []string{filepath.Join(root, "missing-parent", "file"), root} {
		if output, status := runCase(badPath, "cannot-write\n"+end, nil, ""); status == 0 || output == "" {
			t.Fatalf("invalid destination %q: %d %q", badPath, status, output)
		}
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(root, "unwritable")
		mustWriteFile(t, blocked, []byte("old"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runCase(blocked, "cannot-write\n"+end, nil, "")
		if err := os.Chmod(blocked, 0o600); err != nil {
			t.Fatal(err)
		}
		if status == 0 || output == "" {
			t.Fatalf("permission failure: %d %q", status, output)
		}
	}
	if output, status := runCase(existing, "not-written\n"+end, []string{"SHELLMAN_TEE_FAIL=1"}, ""); status != 37 {
		t.Fatalf("writer failure status: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "through-link\n" {
		t.Fatalf("failed writer changed destination: %q %v", got, err)
	}
	if output, status := runCase(existing, "cancel\n\x03", nil, ""); status == 0 {
		t.Fatalf("interrupted input: %d %q", status, output)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\nprintf 'partial' > \"$1\"\nexit 29\n")
	if output, status := runCase(existing, "full content\n"+end, nil, ""); status != 29 {
		t.Fatalf("partial writer failure: %d %q", status, output)
	}
	if got, err := os.ReadFile(existing); err != nil || string(got) != "partial" {
		t.Fatalf("partial writer bytes: %q %v", got, err)
	}
	mustExecutable(t, filepath.Join(bin, "tee"), "#!/bin/sh\nkill -TERM \"$$\"\n")
	if output, status := runCase(existing, "interrupted writer\n"+end, nil, ""); status == 0 {
		t.Fatalf("interrupted writer: %d %q", status, output)
	}
	if err := os.Remove(filepath.Join(bin, "tee")); err != nil {
		t.Fatal(err)
	}
	if output, status := runCase(existing, "missing tee\n"+end, nil, ""); status == 0 {
		t.Fatalf("missing tee: %d %q", status, output)
	}
	if output, status := runCase(existing, "missing tee\n"+end, nil, "printf 'AFTER\\n'\n"); status != 0 || !strings.Contains(output, "AFTER") {
		t.Fatalf("later-command masking or shell continuation: %d %q", status, output)
	}
	// Nonterminal and closed stdin fail before a destination is opened.
	scriptPath := filepath.Join(root, "write.sh")
	mustWriteFile(t, scriptPath, []byte("#!/bin/sh\n"+expandedGeneratedFileWrite54(t, "${pathToFile}")))
	for _, closeFD := range []bool{false, true} {
		var command *exec.Cmd
		if closeFD {
			quoted := "'" + strings.ReplaceAll(scriptPath, "'", "'\\''") + "'"
			command = exec.Command("/bin/sh", "-c", "exec 0<&-; exec /bin/sh "+quoted)
		} else {
			command = exec.Command("/bin/sh", scriptPath)
			command.Stdin = strings.NewReader("redirected\n")
		}
		unopened := filepath.Join(root, "must-not-open")
		command.Env = append(os.Environ(), "pathToFile="+unopened, "PATH="+bin)
		output, err := command.CombinedOutput()
		if exitCode(err) != 2 || !strings.Contains(string(output), "interactive terminal input is required") {
			t.Fatalf("nonterminal or closed input: %d %q", exitCode(err), output)
		}
		if _, err := os.Stat(unopened); !os.IsNotExist(err) {
			t.Fatalf("nonterminal input opened destination: %v", err)
		}
	}
}

func expandedGeneratedFilesIterate55(t *testing.T, pattern, action string) string {
	t.Helper()
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	body, ok := snippets["filesystem.files-iterate"].Body.([]any)
	if !ok || !reflect.DeepEqual(body, currentFilesIterateBody) {
		t.Fatalf("generated entry-55 body = %#v", snippets["filesystem.files-iterate"].Body)
	}
	lines := make([]string, len(body))
	for i, value := range body {
		line, ok := value.(string)
		if !ok {
			t.Fatalf("body line %d = %#v", i, value)
		}
		lines[i] = line
	}
	joined := strings.Join(lines, "\n")
	patternPlaceholder := `${1|*,*.{jpg\,png}|}`
	actionPlaceholder := `${2:printf '%s\n' "\${file\}"}`
	if strings.Count(joined, patternPlaceholder) != 1 || strings.Count(joined, actionPlaceholder) != 1 {
		t.Fatalf("entry-55 placeholder contract = %q", joined)
	}
	joined = strings.Replace(joined, patternPlaceholder, pattern, 1)
	joined = strings.Replace(joined, actionPlaceholder, action, 1)
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(joined, `\,`, `,`), `\$`, `$`), `\}`, `}`)
}

func TestFilesIterate55GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	requireCommand(t, "bash")
	requireCommand(t, "shellcheck")
	defaultAction := `printf '%s\n' "${file}"`
	for _, pattern := range []string{`*`, `*.{jpg,png}`} {
		body := expandedGeneratedFilesIterate55(t, pattern, defaultAction)
		script := filepath.Join(t.TempDir(), "lint.sh")
		mustWriteFile(t, script, []byte("#!/bin/bash\n"+body))
		run(t, ".", "bash", "-n", script)
		if output, err := runCommand(".", nil, "shellcheck", "-s", "bash", script); err != nil {
			t.Fatalf("ShellCheck %q: %v %q", pattern, err, output)
		}
		if pattern == `*` {
			run(t, ".", "sh", "-n", script)
			if output, err := runCommand(".", nil, "shellcheck", "-s", "sh", script); err != nil {
				t.Fatalf("POSIX ShellCheck: %v %q", err, output)
			}
		}
	}
	root := t.TempDir()
	work := filepath.Join(root, "work")
	mustMkdirAll(t, work)
	script := filepath.Join(root, "iterate.sh")
	for index, tc := range []struct{ name, pattern, shell string }{
		{"-n", `*`, "/bin/sh"},
		{"-e", `*`, "/bin/sh"},
		{"back\\slash", `*`, "/bin/sh"},
		{"space name", `*`, "/bin/sh"},
		{"tab\tname", `*`, "/bin/sh"},
		{"line\nbreak", `*`, "/bin/sh"},
		{"glob*?name", `*`, "/bin/sh"},
		{"dollar$sign", `*`, "/bin/sh"},
		{"single'quote", `*`, "/bin/sh"},
		{"double\"quote", `*`, "/bin/sh"},
		{"-n.jpg", `*.{jpg,png}`, "/bin/bash"},
		{"-e.png", `*.{jpg,png}`, "/bin/bash"},
	} {
		single := filepath.Join(root, fmt.Sprintf("single-%d", index))
		mustMkdirAll(t, single)
		mustWriteFile(t, filepath.Join(single, tc.name), []byte("x"))
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+expandedGeneratedFilesIterate55(t, tc.pattern, defaultAction)))
		output, err := runCommand(single, nil, tc.shell, script)
		if err != nil || !bytes.Equal(output, []byte(tc.name+"\n")) {
			t.Fatalf("default printf %q (%s): %v %q", tc.name, tc.shell, err, output)
		}
	}
	runBody := func(pattern, action, suffix, shell string) ([]byte, int) {
		t.Helper()
		mustWriteFile(t, script, []byte("#!/bin/sh\n"+expandedGeneratedFilesIterate55(t, pattern, action)+suffix))
		output, err := runCommand(work, nil, shell, script)
		return output, exitCode(err)
	}
	if output, status := runBody(`*`, `printf '%s\n' "$file"`, "", "/bin/sh"); status != 0 || len(output) != 0 {
		t.Fatalf("empty directory: %d %q", status, output)
	}
	for _, name := range []string{"ordinary", ".hidden", "space name", "tab\tname", "line\nbreak", "glob*?name", "-leading", "single'quote", "double\"quote", "back\\slash", "dollar$sign", "semi;colon", "round(parens)", "square[brackets]", "$(touch MARKER)", "photo.jpg", "photo.png"} {
		mustWriteFile(t, filepath.Join(work, name), []byte("x"))
	}
	mustMkdirAll(t, filepath.Join(work, "directory"))
	mustWriteFile(t, filepath.Join(work, "directory", "nested"), []byte("x"))
	if err := os.Symlink("ordinary", filepath.Join(work, "link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("absent", filepath.Join(work, "broken-link")); err != nil {
		t.Fatal(err)
	}
	output, status := runBody(`*`, `printf '%s\000' "$file"`, "", "/bin/sh")
	if status != 0 {
		t.Fatalf("mixed entries: %d %q", status, output)
	}
	got := strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00")
	sort.Strings(got)
	want := []string{"$(touch MARKER)", "-leading", "back\\slash", "broken-link", "directory", "dollar$sign", "double\"quote", "glob*?name", "line\nbreak", "link", "ordinary", "photo.jpg", "photo.png", "round(parens)", "semi;colon", "single'quote", "space name", "square[brackets]", "tab\tname"}
	sort.Strings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("immediate entries = %#v, want %#v", got, want)
	}
	if _, err := os.Stat(filepath.Join(work, "MARKER")); !os.IsNotExist(err) {
		t.Fatalf("filename executed as shell syntax: %v", err)
	}
	output, status = runBody(`*.{jpg,png}`, `printf '%s\000' "$file"`, "", "/bin/bash")
	if status != 0 {
		t.Fatalf("brace glob: %d %q", status, output)
	}
	got = strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00")
	sort.Strings(got)
	if !reflect.DeepEqual(got, []string{"photo.jpg", "photo.png"}) {
		t.Fatalf("brace glob matches = %#v", got)
	}
	if output, status := runBody(`missing/*`, defaultAction, "", "/bin/sh"); status != 0 || len(output) != 0 {
		t.Fatalf("missing directory pattern: %d %q", status, output)
	}
	if os.Geteuid() != 0 {
		blocked := filepath.Join(work, "blocked")
		mustMkdirAll(t, blocked)
		mustWriteFile(t, filepath.Join(blocked, "child"), []byte("x"))
		if err := os.Chmod(blocked, 0); err != nil {
			t.Fatal(err)
		}
		output, status := runBody(`blocked/*`, defaultAction, "", "/bin/sh")
		if err := os.Chmod(blocked, 0o700); err != nil {
			t.Fatal(err)
		}
		if status != 0 || len(output) != 0 {
			t.Fatalf("unreadable directory pattern: %d %q", status, output)
		}
	}
	output, status = runBody(`ordinary photo.jpg`, `case "$file" in ordinary) (exit 23);; *) printf '%s\n' "$file";; esac`, "", "/bin/sh")
	if status != 23 || string(output) != "photo.jpg\n" {
		t.Fatalf("first action failure and later success: %d %q", status, output)
	}
	output, status = runBody(`ordinary photo.jpg`, `case "$file" in ordinary) (exit 23);; *) printf '%s\n' "$file";; esac`, "true\n", "/bin/sh")
	if status != 0 || string(output) != "photo.jpg\n" {
		t.Fatalf("later-command masking: %d %q", status, output)
	}
	output, status = runBody(`ordinary`, `seen=yes; cd ..`, "printf '%s:%s\\n' \"$seen\" \"$PWD\"\n", "/bin/sh")
	if status != 0 || string(output) != "yes:"+root+"\n" {
		t.Fatalf("action shell and cwd effect: %d %q", status, output)
	}
	marker := filepath.Join(work, "INJECTED")
	injection := `$(touch ` + marker + `)`
	if output, status := runBody(injection, defaultAction, "", "/bin/sh"); status != 0 || len(output) != 0 {
		t.Fatalf("editable expression behavior: %d %q", status, output)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Fatalf("editable shell expression did not execute as expected: %v", err)
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

func TestFindFilesOrDirectories56GeneratedBehavior(t *testing.T) {
	requireCommand(t, "bash")
	requireCommand(t, "find")
	requireCommand(t, "printf")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.find-files-or-directories"]
	if !reflect.DeepEqual(item.Prefix, []any{"find file", "find directory"}) || item.Description != "find files (-type f) or directories (-type d) by name or pattern (*.jpg)" || item.Body != currentFindFilesOrDirectoriesBody {
		t.Fatalf("entry-56 contract = %#v", item)
	}
	body := item.Body.(string)
	if strings.Contains(body, "${1:") || strings.Contains(body, "${1|") {
		t.Fatal("unexpected VS Code placeholder")
	}
	quote := func(s string) string { return "'" + strings.ReplaceAll(s, "'", "'\\''") + "'" }
	root := t.TempDir()
	for _, name := range []string{"plain.txt", ".hidden", "space name", "tab\tname", "line\nname", "glob[*]", "-leading", "-", "$(touch injected)"} {
		mustWriteFile(t, filepath.Join(root, name), nil)
	}
	for _, name := range []string{"nested", "empty"} {
		if err := os.Mkdir(filepath.Join(root, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	mustWriteFile(t, filepath.Join(root, "nested", "deep.txt"), nil)
	if err := os.Symlink("plain.txt", filepath.Join(root, "valid-link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("missing", filepath.Join(root, "broken-link")); err != nil {
		t.Fatal(err)
	}
	run := func(searchRoot, pattern, kind, pathEnv, workdir string) ([]string, int, string) {
		t.Helper()
		script := strings.Replace(body, `_shellman_find_root="/path/to/search"`, `_shellman_find_root=`+quote(searchRoot), 1)
		script = strings.Replace(script, `_shellman_find_pattern="criteria"`, `_shellman_find_pattern=`+quote(pattern), 1)
		if kind == "d" {
			script = strings.Replace(script, "-type f -name", "-type d -name", 1)
		}
		script += "_shellman_find_status=$?\nfor _shellman_path in \"${resultArray[@]}\"; do printf '%s\\0' \"$_shellman_path\"; done\nexit \"$_shellman_find_status\"\n"
		cmd := exec.Command("bash", "-c", script)
		if workdir != "" {
			cmd.Dir = workdir
		}
		if pathEnv != "" {
			cmd.Env = append(os.Environ(), "PATH="+pathEnv)
		}
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		output, err := cmd.Output()
		if len(output) > 0 && output[len(output)-1] != 0 {
			t.Fatalf("unterminated output: %q", output)
		}
		var names []string
		for _, part := range bytes.Split(bytes.TrimSuffix(output, []byte{0}), []byte{0}) {
			if len(part) > 0 {
				names = append(names, string(part))
			}
		}
		sort.Strings(names)
		return names, exitCode(err), stderr.String()
	}
	check := func(label string, got []string, status int, want []string) {
		t.Helper()
		sort.Strings(want)
		if status != 0 || !reflect.DeepEqual(got, want) {
			t.Fatalf("%s: status=%d got=%q want=%q", label, status, got, want)
		}
	}
	files := []string{"plain.txt", ".hidden", "space name", "tab\tname", "line\nname", "glob[*]", "-leading", "-", "$(touch injected)"}
	wantFiles := make([]string, len(files))
	for i, n := range files {
		wantFiles[i] = filepath.Join(root, n)
	}
	got, status, stderr := run(root, "*", "f", "", "")
	check("immediate regular files", got, status, wantFiles)
	got, status, _ = run(root, "*.txt", "f", "", "")
	check("name glob", got, status, []string{filepath.Join(root, "plain.txt")})
	got, status, _ = run(root, "$(touch injected)", "f", "", "")
	check("inert pattern", got, status, []string{filepath.Join(root, "$(touch injected)")})
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("pattern executed: %v", err)
	}
	inertRoot := filepath.Join(t.TempDir(), "$(touch root-injected)")
	if err := os.Mkdir(inertRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(inertRoot, "plain"), nil)
	got, status, _ = run(inertRoot, "*", "f", "", "")
	check("inert root", got, status, []string{filepath.Join(inertRoot, "plain")})
	if _, err := os.Stat(filepath.Join(inertRoot, "root-injected")); !os.IsNotExist(err) {
		t.Fatalf("root executed: %v", err)
	}
	patternRoot := filepath.Join(t.TempDir(), "root[*]?\\")
	if err := os.Mkdir(patternRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(patternRoot, "plain"), nil)
	got, status, _ = run(patternRoot, "*", "f", "", "")
	check("glob characters in root", got, status, []string{filepath.Join(patternRoot, "plain")})
	got, status, _ = run(filepath.Join(root, "plain.txt"), "*.txt", "f", "", "")
	check("regular file root", got, status, []string{filepath.Join(root, "plain.txt")})
	got, status, _ = run(root, "*", "d", "", "")
	check("immediate directories", got, status, []string{root, filepath.Join(root, "nested"), filepath.Join(root, "empty")})
	got, status, _ = run(filepath.Join(root, "empty"), "*", "f", "", "")
	check("empty directory", got, status, nil)
	parent := t.TempDir()
	if err := os.Mkdir(filepath.Join(parent, "-root"), 0o755); err != nil {
		t.Fatal(err)
	}
	mustWriteFile(t, filepath.Join(parent, "-root", "-"), nil)
	got, status, _ = run("-root", "*", "f", "", parent)
	check("leading hyphen root", got, status, []string{"./-root/-"})
	for _, bad := range []string{"", filepath.Join(root, "absent")} {
		got, status, stderr := run(bad, "*", "f", "", "")
		if status == 0 || len(got) != 0 || stderr == "" {
			t.Fatalf("bad root %q: status=%d names=%q stderr=%q", bad, status, got, stderr)
		}
	}
	missingBin := t.TempDir()
	got, status, stderr = run(root, "*", "f", missingBin, "")
	if status != 127 || len(got) != 0 || !strings.Contains(stderr, "find") {
		t.Fatalf("missing find: status=%d names=%q stderr=%q", status, got, stderr)
	}
	fakeBin := t.TempDir()
	mustExecutable(t, filepath.Join(fakeBin, "find"), "#!/bin/sh\nprintf 'partial\\0'\nexit 23\n")
	got, status, stderr = run(root, "*", "f", fakeBin, "")
	if status != 23 || !reflect.DeepEqual(got, []string{"partial"}) {
		t.Fatalf("partial find failure: status=%d names=%q stderr=%q", status, got, stderr)
	}
	unreadable := filepath.Join(root, "unreadable")
	if err := os.Mkdir(unreadable, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(unreadable, 0o755)
	got, status, stderr = run(root, "*", "d", "", "")
	// A depth-one search prunes unreadable children, so it can still return them by name.
	check("unreadable child pruned", got, status, []string{root, filepath.Join(root, "nested"), filepath.Join(root, "empty"), unreadable})
	got, status, stderr = run(unreadable, "*", "f", "", "")
	if status == 0 || len(got) != 0 || stderr == "" {
		t.Fatalf("unreadable root: status=%d names=%q stderr=%q", status, got, stderr)
	}
}

func TestIfDirectoryExists57GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-directory-exists"]
	if item.Prefix != "if directory exists" || item.Description != "check if a directory exists" || !reflect.DeepEqual(item.Body, v6IfDirectoryExistsBody) {
		t.Fatalf("entry-57 contract: %#v", item)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/directory,${pathToDirectory}|}`
	actionPlaceholder := `${2:echo "directory \"${1}\" exists"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(choice, action string) string {
		t.Helper()
		s := strings.Replace(body, pathPlaceholder, choice, 1)
		return strings.Replace(s, actionPlaceholder, action, 1)
	}
	defaultAction := `echo "directory \"${pathToDirectory}\" exists"`
	defaultExpanded := expand(`${pathToDirectory}`, defaultAction)
	if defaultExpanded != "if [ -d \"${pathToDirectory}\" ]; then\n\t"+defaultAction+"\nfi\n" {
		t.Fatalf("default generated expansion: %q", defaultExpanded)
	}
	root := t.TempDir()
	mustWriteFile(t, filepath.Join(root, "file"), nil)
	if err := os.Mkdir(filepath.Join(root, "directory"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("directory", filepath.Join(root, "directory-link")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("absent", filepath.Join(root, "broken-link")); err != nil {
		t.Fatal(err)
	}
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		if err := os.Mkdir(filepath.Join(root, name), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	unreadable := filepath.Join(root, "unreadable")
	if err := os.Mkdir(unreadable, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(unreadable, 0o755)
	blocked := filepath.Join(root, "blocked")
	if err := os.Mkdir(blocked, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(blocked, "child"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(blocked, 0o000); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(blocked, 0o755)
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${pathToDirectory}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "pathToDirectory="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	printTrue := `printf 'TRUE\n'`
	check := func(label, path string, want bool, cwd string) {
		t.Helper()
		out, status, stderr := run(path, printTrue, "", false, cwd)
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		if status != 0 || out != wantOut || stderr != "" {
			t.Fatalf("%s: status=%d stdout=%q stderr=%q want=%q", label, status, out, stderr, wantOut)
		}
		out, status, stderr = run(path, printTrue, "", true, cwd)
		if status != 0 || out != wantOut || stderr != "" {
			t.Fatalf("%s set -e: status=%d stdout=%q stderr=%q want=%q", label, status, out, stderr, wantOut)
		}
	}
	check("existing absolute directory", filepath.Join(root, "directory"), true, "")
	check("existing relative directory", "directory", true, root)
	check("missing", filepath.Join(root, "missing"), false, "")
	check("regular file", filepath.Join(root, "file"), false, "")
	check("directory symlink", filepath.Join(root, "directory-link"), true, "")
	check("broken symlink", filepath.Join(root, "broken-link"), false, "")
	check("empty path", "", false, "")
	check("unreadable directory", unreadable, true, "")
	check("non-searchable directory itself", blocked, true, "")
	if os.Geteuid() != 0 {
		check("child under non-searchable ancestor", filepath.Join(blocked, "child"), false, "")
	}
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), true, "")
	}
	check("relative leading hyphen", "-leading", true, root)
	check("relative literal hyphen", "-", true, root)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(filepath.Join(root, "directory"), defaultAction, "", false, "")
	wantDefault := "directory \"" + filepath.Join(root, "directory") + "\" exists\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	out, status, stderr = run(filepath.Join(root, "missing"), defaultAction, "", false, "")
	if status != 0 || out != "STATUS:0\n" || stderr != "" {
		t.Fatalf("false default: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(filepath.Join(root, "directory"), "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e true branch failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("true branch failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(filepath.Join(root, "missing"), "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch skip: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(filepath.Join(root, "directory"), `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
	if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("caller shell scope: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
}

func TestIfFileExecutable58GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file-executable"]
	if item.Prefix != "if file executable" || item.Description != "check if file is executable" || !reflect.DeepEqual(item.Body, currentIfFileExecutableBody) {
		t.Fatalf("entry-58 contract: %#v", item)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/file,${filePath}|}`
	actionPlaceholder := `${2:echo "file \"${1}\" is executable"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(choice, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(body, pathPlaceholder, choice, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "file \"${filePath}\" is executable"`
	defaultExpanded := expand(`${filePath}`, defaultAction)
	wantDefaultExpanded := "if (\n\t_shellman_executable_path=\"${filePath}\"\n\t[ -f \"$_shellman_executable_path\" ] && [ -x \"$_shellman_executable_path\" ]\n); then\n\t" + defaultAction + "\nfi\n"
	if defaultExpanded != wantDefaultExpanded {
		t.Fatalf("default expansion = %q", defaultExpanded)
	}
	root := t.TempDir()
	file := func(name string, mode os.FileMode) string {
		t.Helper()
		p := filepath.Join(root, name)
		if err := os.WriteFile(p, []byte("not a runnable program\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(p, mode); err != nil {
			t.Fatal(err)
		}
		return p
	}
	execFile := file("executable", 0o755)
	noExecFile := file("not-executable", 0o644)
	execDir := filepath.Join(root, "executable-directory")
	if err := os.Mkdir(execDir, 0o755); err != nil {
		t.Fatal(err)
	}
	noSearchDir := filepath.Join(root, "non-searchable-directory")
	if err := os.Mkdir(noSearchDir, 0o600); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(noSearchDir, 0o755)
	if err := os.Symlink("executable", filepath.Join(root, "link-executable")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("not-executable", filepath.Join(root, "link-not-executable")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink("absent", filepath.Join(root, "broken-link")); err != nil {
		t.Fatal(err)
	}
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		file(name, 0o755)
	}
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	printTrue := `printf 'TRUE\n'`
	check := func(label, path string, want bool, cwd string) {
		t.Helper()
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(path, printTrue, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s set-e=%v: status=%d stdout=%q stderr=%q want=%q", label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	check("executable regular", execFile, true, "")
	check("non-executable regular", noExecFile, false, "")
	check("executable directory is not a file", execDir, false, "")
	check("non-searchable directory is not a file", noSearchDir, false, "")
	check("symlink to executable regular", filepath.Join(root, "link-executable"), true, "")
	check("symlink to non-executable regular", filepath.Join(root, "link-not-executable"), false, "")
	check("broken symlink", filepath.Join(root, "broken-link"), false, "")
	check("missing", filepath.Join(root, "missing"), false, "")
	check("empty path", "", false, "")
	check("relative executable", "executable", true, root)
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), true, "")
	}
	check("relative leading hyphen", "-leading", true, root)
	check("relative literal hyphen", "-", true, root)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(execFile, defaultAction, "", false, "")
	wantDefault := "file \"" + execFile + "\" is executable\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(execFile, "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(noExecFile, "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(execFile, `marker=after; printf '%s:%s\n' "$marker" "$_shellman_executable_path"`, `marker=before; _shellman_executable_path=outer`, false, "")
	if status != 0 || out != "after:outer\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("action scope/subshell isolation: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
	// The 0755 regular fixture is only tested with -f and -x; it is never executed.
}

func TestIfFileSymbolicLink59GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file-exists-and-is-symbolic-link"]
	historicalBody := []any{
		`if [ -h "${1|/path/to/file,${filePath}|}" ]; then`,
		"\t${2:echo \"Path \\\"${1}\\\" is a symbolic link\"}",
		"fi\n",
	}
	if item.Prefix != "if file link" || item.Description != "if given path is a symbolic link" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-59 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file-exists-and-is-symbolic-link.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/file,${filePath}|}`
	actionPlaceholder := `${2:echo "Path \"${1}\" is a symbolic link"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(choice, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(body, pathPlaceholder, choice, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "Path \"${filePath}\" is a symbolic link"`
	if got, want := expand(`${filePath}`, defaultAction), "if [ -h \"${filePath}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("default expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	file := filepath.Join(root, "regular")
	if err := os.WriteFile(file, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(root, "directory")
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	link := func(name, target string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		return path
	}
	validFile := link("valid-file", file)
	validDirectory := link("valid-directory", directory)
	broken := link("broken", filepath.Join(root, "absent"))
	danglingRelative := link("dangling-relative", "absent-relative")
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		link(name, file)
	}
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	check := func(label, path string, want bool, cwd string) {
		t.Helper()
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(path, `printf 'TRUE\n'`, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s set-e=%v: status=%d stdout=%q stderr=%q want=%q", label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	check("valid file link", validFile, true, "")
	check("valid directory link", validDirectory, true, "")
	check("broken absolute link", broken, true, "")
	check("dangling relative link", danglingRelative, true, "")
	check("regular file", file, false, "")
	check("directory", directory, false, "")
	check("missing", filepath.Join(root, "missing"), false, "")
	check("empty", "", false, "")
	check("relative link", "valid-file", true, root)
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), true, "")
	}
	check("relative leading hyphen", "-leading", true, root)
	check("relative literal hyphen", "-", true, root)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(validFile, defaultAction, "", false, "")
	wantDefault := "Path \"" + validFile + "\" is a symbolic link\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(validFile, "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(file, "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(validFile, `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
	if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("caller shell scope: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
}

func TestIfFileExists60GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file-exists"]
	historicalBody := []any{
		`if [ -f "${1|/path/to/file,${filePath}|}" ]; then`,
		"\t${2:echo \"File \\\"${1}\\\" exists\"}",
		"fi\n",
	}
	if item.Prefix != "if file exists" || item.Description != "check if the file exists" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-60 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file-exists.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/file,${filePath}|}`
	actionPlaceholder := `${2:echo "File \"${1}\" exists"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(choice, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(body, pathPlaceholder, choice, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "File \"${filePath}\" exists"`
	if got, want := expand(`${filePath}`, defaultAction), "if [ -f \"${filePath}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("default expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	file := func(name string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	regular := file("regular")
	directory := filepath.Join(root, "directory")
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	link := func(name, target string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		return path
	}
	fileLink := link("file-link", regular)
	directoryLink := link("directory-link", directory)
	brokenLink := link("broken-link", filepath.Join(root, "absent"))
	relativeLink := link("relative-link", "regular")
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		file(name)
	}
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	check := func(label, path string, want bool, cwd string) {
		t.Helper()
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(path, `printf 'TRUE\n'`, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s set-e=%v: status=%d stdout=%q stderr=%q want=%q", label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	check("regular file", regular, true, "")
	check("directory", directory, false, "")
	check("file symlink", fileLink, true, "")
	check("directory symlink", directoryLink, false, "")
	check("broken symlink", brokenLink, false, "")
	check("relative target symlink", relativeLink, true, "")
	check("missing", filepath.Join(root, "missing"), false, "")
	check("empty", "", false, "")
	check("relative regular", "regular", true, root)
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), true, "")
	}
	check("relative leading hyphen", "-leading", true, root)
	check("relative literal hyphen", "-", true, root)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(regular, defaultAction, "", false, "")
	wantDefault := "File \"" + regular + "\" exists\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(regular, "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(directory, "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(regular, `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
	if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("caller shell scope: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
}

func TestIfFileNotEmpty61GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file-not-empty"]
	historicalBody := []any{
		`if [ -s "${1|/path/to/file,${filePath}|}" ]; then`,
		"\t${2:echo \"File \\\"${1}\\\" is not empty\"}",
		"fi\n",
	}
	if item.Prefix != "if file not empty" || item.Description != "check if file size is greater than zero" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-61 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file-not-empty.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/file,${filePath}|}`
	actionPlaceholder := `${2:echo "File \"${1}\" is not empty"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(pathChoice, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(body, pathPlaceholder, pathChoice, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "File \"${filePath}\" is not empty"`
	if got, want := expand(`${filePath}`, defaultAction), "if [ -s \"${filePath}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("default expansion = %q, want %q", got, want)
	}
	if got, want := expand(`/path/to/file`, `echo "File \"/path/to/file\" is not empty"`), "if [ -s \"/path/to/file\" ]; then\n\techo \"File \\\"/path/to/file\\\" is not empty\"\nfi\n"; got != want {
		t.Fatalf("literal-choice expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	file := func(name string, contents []byte) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, contents, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	empty := file("empty", nil)
	one := file("one", []byte("x"))
	nonempty := file("nonempty", []byte("contents"))
	sparse := file("sparse", nil)
	if err := os.Truncate(sparse, 1<<20); err != nil {
		t.Fatal(err)
	}
	emptyDir := filepath.Join(root, "empty-dir")
	nonemptyDir := filepath.Join(root, "nonempty-dir")
	for _, dir := range []string{emptyDir, nonemptyDir} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	file("nonempty-dir/child", nil)
	link := func(name, target string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		return path
	}
	emptyLink := link("empty-link", empty)
	nonemptyLink := link("nonempty-link", nonempty)
	brokenLink := link("broken-link", filepath.Join(root, "absent"))
	unreadable := file("unreadable", []byte("x"))
	if err := os.Chmod(unreadable, 0); err != nil {
		t.Fatal(err)
	}
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		file(name, []byte("x"))
	}
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	check := func(label, path string, want bool, cwd string) {
		t.Helper()
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(path, `printf 'TRUE\n'`, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s set-e=%v: status=%d stdout=%q stderr=%q want=%q", label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	check("zero-byte file", empty, false, "")
	check("one-byte file", one, true, "")
	check("nonempty file", nonempty, true, "")
	check("sparse file", sparse, true, "")
	for _, dir := range []string{emptyDir, nonemptyDir} {
		info, err := os.Stat(dir)
		if err != nil {
			t.Fatal(err)
		}
		check("directory "+dir, dir, info.Size() > 0, "")
	}
	check("symlink to empty file", emptyLink, false, "")
	check("symlink to nonempty file", nonemptyLink, true, "")
	check("broken symlink", brokenLink, false, "")
	check("unreadable regular file", unreadable, true, "")
	check("missing", filepath.Join(root, "missing"), false, "")
	check("empty path", "", false, "")
	check("relative regular file", "one", true, root)
	check("relative symlink", "nonempty-link", true, root)
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), true, "")
	}
	check("relative leading hyphen", "-leading", true, root)
	check("relative literal hyphen", "-", true, root)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(one, defaultAction, "", false, "")
	wantDefault := "File \"" + one + "\" is not empty\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(one, "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(empty, "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(one, `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
	if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("caller shell scope: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
}

func TestIfFileReadable62GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file-readable"]
	historicalBody := []any{
		`if [ -r "${1|/path/to/file,${filePath}|}" ]; then`,
		"\t${2:echo \"File \\\"${1}\\\" is readable\"}",
		"fi\n",
	}
	if item.Prefix != "if file readable" || item.Description != "check if file readable" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-62 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file-readable.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/file,${filePath}|}`
	actionPlaceholder := `${2:echo "File \"${1}\" is readable"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(pathChoice, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(body, pathPlaceholder, pathChoice, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "File \"${filePath}\" is readable"`
	if got, want := expand(`${filePath}`, defaultAction), "if [ -r \"${filePath}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("variable-choice expansion = %q, want %q", got, want)
	}
	if got, want := expand(`/path/to/file`, `echo "File \"/path/to/file\" is readable"`), "if [ -r \"/path/to/file\" ]; then\n\techo \"File \\\"/path/to/file\\\" is readable\"\nfi\n"; got != want {
		t.Fatalf("literal-choice expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	file := func(name string, contents []byte, mode os.FileMode) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, contents, mode); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(path, mode); err != nil {
			t.Fatal(err)
		}
		return path
	}
	empty := file("empty-readable", nil, 0o600)
	nonempty := file("nonempty-readable", []byte("contents"), 0o600)
	unreadable := file("mode-000-file", []byte("contents"), 0)
	readableDir := filepath.Join(root, "readable-dir")
	unreadableDir := filepath.Join(root, "mode-000-dir")
	for _, dir := range []string{readableDir, unreadableDir} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chmod(unreadableDir, 0); err != nil {
		t.Fatal(err)
	}
	link := func(name, target string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		return path
	}
	readableLink := link("readable-link", nonempty)
	unreadableLink := link("mode-000-link", unreadable)
	brokenLink := link("broken-link", filepath.Join(root, "missing-target"))
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		file(name, []byte("x"), 0o600)
	}
	probeReadable := func(path, cwd string) bool {
		t.Helper()
		cmd := exec.Command("sh", "-c", `[ -r "$filePath" ]`)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		err := cmd.Run()
		if err == nil {
			return true
		}
		if exitCode(err) != 1 {
			t.Fatalf("readability probe for %q: %v", path, err)
		}
		return false
	}
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	check := func(label, path, cwd string, want bool) {
		t.Helper()
		if measured := probeReadable(path, cwd); measured != want {
			t.Fatalf("%s: measured -r=%v, want %v", label, measured, want)
		}
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(path, `printf 'TRUE\n'`, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s set-e=%v: status=%d stdout=%q stderr=%q want=%q", label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	check("empty readable file", empty, "", true)
	check("nonempty readable file", nonempty, "", true)
	check("readable directory", readableDir, "", true)
	for _, fixture := range []struct{ label, path string }{
		{"mode-000 regular file", unreadable},
		{"mode-000 directory", unreadableDir},
		{"symlink to mode-000 file", unreadableLink},
	} {
		measured := probeReadable(fixture.path, "")
		t.Logf("effective UID %d: %s has -r=%v", os.Geteuid(), fixture.label, measured)
		check(fixture.label, fixture.path, "", measured)
	}
	check("symlink to readable file", readableLink, "", true)
	check("broken symlink", brokenLink, "", false)
	check("missing", filepath.Join(root, "missing"), "", false)
	check("empty path", "", "", false)
	check("relative file", "nonempty-readable", root, true)
	check("relative symlink", "readable-link", root, true)
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), "", true)
	}
	check("relative leading hyphen", "-leading", root, true)
	check("relative literal hyphen", "-", root, true)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(nonempty, defaultAction, "", false, "")
	wantDefault := "File \"" + nonempty + "\" is readable\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(nonempty, "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(brokenLink, "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(nonempty, `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
	if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("caller shell scope: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
}

func TestIfFileWriteable63GeneratedBehavior(t *testing.T) {
	requireCommand(t, "sh")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file-writeable"]
	historicalBody := []any{
		`if [ -w "${1|/path/to/file,${filePath}|}" ]; then`,
		"\t${2:echo \"File \\\"${1}\\\" is writeable\"}",
		"fi\n",
	}
	if item.Prefix != "if file writeable" || item.Description != "if file writeable" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-63 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file-writeable.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/file,${filePath}|}`
	actionPlaceholder := `${2:echo "File \"${1}\" is writeable"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling/order: %q", body)
	}
	expand := func(pathChoice, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(body, pathPlaceholder, pathChoice, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "File \"${filePath}\" is writeable"`
	if got, want := expand(`${filePath}`, defaultAction), "if [ -w \"${filePath}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("variable-choice expansion = %q, want %q", got, want)
	}
	if got, want := expand(`/path/to/file`, `echo "File \"/path/to/file\" is writeable"`), "if [ -w \"/path/to/file\" ]; then\n\techo \"File \\\"/path/to/file\\\" is writeable\"\nfi\n"; got != want {
		t.Fatalf("literal-choice expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	file := func(name string, contents []byte, mode os.FileMode) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, contents, mode); err != nil {
			t.Fatal(err)
		}
		if err := os.Chmod(path, mode); err != nil {
			t.Fatal(err)
		}
		return path
	}
	empty := file("empty-writeable", nil, 0o600)
	nonempty := file("nonempty-writeable", []byte("contents"), 0o600)
	unwriteable := file("mode-000-file", []byte("contents"), 0)
	writeableDir := filepath.Join(root, "writeable-dir")
	unwriteableDir := filepath.Join(root, "mode-000-dir")
	for _, dir := range []string{writeableDir, unwriteableDir} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.Chmod(unwriteableDir, 0); err != nil {
		t.Fatal(err)
	}
	link := func(name, target string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		return path
	}
	writeableLink := link("writeable-link", nonempty)
	writeableDirLink := link("writeable-dir-link", writeableDir)
	unwriteableLink := link("mode-000-link", unwriteable)
	brokenLink := link("broken-link", filepath.Join(root, "missing-target"))
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		file(name, []byte("x"), 0o600)
	}
	probeWriteable := func(path, cwd string) bool {
		t.Helper()
		cmd := exec.Command("sh", "-c", `[ -w "$filePath" ]`)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		err := cmd.Run()
		if err == nil {
			return true
		}
		if exitCode(err) != 1 {
			t.Fatalf("writability probe for %q: %v", path, err)
		}
		return false
	}
	run := func(path, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command("sh", "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = append(os.Environ(), "filePath="+path)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	check := func(label, path, cwd string, want bool) {
		t.Helper()
		if measured := probeWriteable(path, cwd); measured != want {
			t.Fatalf("%s: measured -w=%v, want %v", label, measured, want)
		}
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(path, `printf 'TRUE\n'`, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s set-e=%v: status=%d stdout=%q stderr=%q want=%q", label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	check("empty writeable file", empty, "", true)
	check("nonempty writeable file", nonempty, "", true)
	check("writeable directory", writeableDir, "", true)
	for _, fixture := range []struct{ label, path string }{
		{"mode-000 regular file", unwriteable},
		{"mode-000 directory", unwriteableDir},
		{"symlink to mode-000 file", unwriteableLink},
	} {
		measured := probeWriteable(fixture.path, "")
		t.Logf("effective UID %d: %s has -w=%v", os.Geteuid(), fixture.label, measured)
		check(fixture.label, fixture.path, "", measured)
	}
	check("symlink to writeable file", writeableLink, "", true)
	check("symlink to writeable directory", writeableDirLink, "", true)
	check("broken symlink", brokenLink, "", false)
	check("missing", filepath.Join(root, "missing"), "", false)
	check("empty path", "", "", false)
	check("relative file", "nonempty-writeable", root, true)
	check("relative symlink", "writeable-link", root, true)
	for _, name := range special {
		check("special "+name, filepath.Join(root, name), "", true)
	}
	check("relative leading hyphen", "-leading", root, true)
	check("relative literal hyphen", "-", root, true)
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
	out, status, stderr := run(nonempty, defaultAction, "", false, "")
	wantDefault := "File \"" + nonempty + "\" is writeable\nSTATUS:0\n"
	if status != 0 || out != wantDefault || stderr != "" {
		t.Fatalf("default action: status=%d stdout=%q stderr=%q want=%q", status, out, stderr, wantDefault)
	}
	for _, errexit := range []bool{false, true} {
		out, status, stderr = run(nonempty, "false", "", errexit, "")
		if errexit {
			if status != 1 || out != "" || stderr != "" {
				t.Fatalf("set -e action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
			}
		} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
			t.Fatalf("action failure: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
		out, status, stderr = run(brokenLink, "false", "", errexit, "")
		if status != 0 || out != "STATUS:0\n" || stderr != "" {
			t.Fatalf("false branch: status=%d stdout=%q stderr=%q", status, out, stderr)
		}
	}
	out, status, stderr = run(nonempty, `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
	if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
		t.Fatalf("caller shell scope: status=%d stdout=%q stderr=%q", status, out, stderr)
	}
}

func TestIfFile1NewerThanFile2_64GeneratedBehavior(t *testing.T) {
	requireCommand(t, "bash")
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file1-newer-than-file2"]
	historicalBody := []any{
		`if [ "${1|/path/to/file1,${filePath1}|}" -nt "${2|/path/to/file2,${filePath2}|}" ]; then`,
		"\t${3:echo \"Path \\\"${1}\\\" is newer than path \\\"${2}\\\"\"}",
		"fi\n",
	}
	if item.Prefix != "if file newer" || item.Description != "check if file1 is newer than file2" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-64 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file1-newer-than-file2.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	firstPlaceholder := `${1|/path/to/file1,${filePath1}|}`
	secondPlaceholder := `${2|/path/to/file2,${filePath2}|}`
	actionPlaceholder := `${3:echo "Path \"${1}\" is newer than path \"${2}\""}`
	for _, placeholder := range []string{firstPlaceholder, secondPlaceholder, actionPlaceholder} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q spelling/order changed in %q", placeholder, body)
		}
	}
	expand := func(first, second, action string) string {
		t.Helper()
		return strings.Replace(strings.Replace(strings.Replace(body, firstPlaceholder, first, 1), secondPlaceholder, second, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "Path \"${filePath1}\" is newer than path \"${filePath2}\""`
	if got, want := expand(`${filePath1}`, `${filePath2}`, defaultAction), "if [ \"${filePath1}\" -nt \"${filePath2}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("variable-choice expansion = %q, want %q", got, want)
	}
	literalAction := `echo "Path \"/path/to/file1\" is newer than path \"/path/to/file2\""`
	if got, want := expand(`/path/to/file1`, `/path/to/file2`, literalAction), "if [ \"/path/to/file1\" -nt \"/path/to/file2\" ]; then\n\t"+literalAction+"\nfi\n"; got != want {
		t.Fatalf("literal-choice expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	file := func(name string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		return path
	}
	base := time.Unix(1600000000, 0)
	setTime := func(path string, timestamp time.Time) {
		t.Helper()
		if err := os.Chtimes(path, timestamp, timestamp); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil {
			t.Fatal(err)
		}
		if !info.ModTime().Equal(timestamp) {
			t.Fatalf("mtime for %q = %s, want %s", path, info.ModTime(), timestamp)
		}
	}
	resolutionFile := file("resolution")
	setTime(resolutionFile, base)
	resolution := time.Duration(0)
	for _, delta := range []time.Duration{time.Nanosecond, time.Microsecond, time.Millisecond, time.Second, 2 * time.Second, 10 * time.Second} {
		if err := os.Chtimes(resolutionFile, base.Add(delta), base.Add(delta)); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(resolutionFile)
		if err != nil {
			t.Fatal(err)
		}
		if info.ModTime().After(base) {
			resolution = info.ModTime().Sub(base)
			break
		}
	}
	if resolution == 0 {
		t.Fatal("filesystem could not distinguish a ten-second mtime change")
	}
	t.Logf("observed mtime resolution at fixture timestamp: %s", resolution)
	gap := 10 * time.Second
	if 2*resolution > gap {
		gap = 2 * resolution
	}
	older := file("older")
	newer := file("newer")
	equalA := file("equal-a")
	equalB := file("equal-b")
	for _, path := range []string{older, equalA, equalB} {
		setTime(path, base)
	}
	setTime(newer, base.Add(gap))
	statOlder, err := os.Stat(older)
	if err != nil {
		t.Fatal(err)
	}
	statNewer, err := os.Stat(newer)
	if err != nil {
		t.Fatal(err)
	}
	statEqualA, err := os.Stat(equalA)
	if err != nil {
		t.Fatal(err)
	}
	statEqualB, err := os.Stat(equalB)
	if err != nil {
		t.Fatal(err)
	}
	if !statNewer.ModTime().After(statOlder.ModTime()) || !statEqualA.ModTime().Equal(statEqualB.ModTime()) {
		t.Fatal("fixture timestamps do not represent strict and equal ordering")
	}
	olderDir := filepath.Join(root, "older-dir")
	newerDir := filepath.Join(root, "newer-dir")
	for _, dir := range []string{olderDir, newerDir} {
		if err := os.Mkdir(dir, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	setTime(olderDir, base)
	setTime(newerDir, base.Add(gap))
	link := func(name, target string) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.Symlink(target, path); err != nil {
			t.Fatal(err)
		}
		return path
	}
	newerLink := link("newer-link", newer)
	olderLink := link("older-link", older)
	brokenLink := link("broken-link", filepath.Join(root, "missing-target"))
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		setTime(file(name), base.Add(gap))
	}
	shells := []string{"bash", "sh"}
	if _, err := exec.LookPath("zsh"); err == nil {
		shells = append(shells, "zsh")
	}
	env := func(first, second string) []string {
		return append(os.Environ(), "filePath1="+first, "filePath2="+second)
	}
	probe := func(shell, expression, first, second, cwd string) int {
		t.Helper()
		cmd := exec.Command(shell, "-c", expression)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = env(first, second)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		status := exitCode(err)
		if status > 1 || stderr.Len() != 0 {
			t.Fatalf("%s probe %q: status=%d stderr=%q", shell, expression, status, stderr.String())
		}
		return status
	}
	run := func(shell, first, second, action, setup string, errexit bool, cwd string) (string, int, string) {
		t.Helper()
		script := setup + "\n"
		if errexit {
			script += "set -e\n"
		}
		script += expand(`${filePath1}`, `${filePath2}`, action)
		script += "snippet_status=$?\nprintf 'STATUS:%s\\n' \"$snippet_status\"\n"
		cmd := exec.Command(shell, "-c", script)
		if cwd != "" {
			cmd.Dir = cwd
		}
		cmd.Env = env(first, second)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	check := func(shell, label, first, second, cwd string, want bool) {
		t.Helper()
		wantProbe := 1
		if want {
			wantProbe = 0
		}
		if got := probe(shell, `[ "$filePath1" -nt "$filePath2" ]`, first, second, cwd); got != wantProbe {
			t.Fatalf("%s %s direct -nt=%d, want %d", shell, label, got, wantProbe)
		}
		wantOut := "STATUS:0\n"
		if want {
			wantOut = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr := run(shell, first, second, `printf 'TRUE\n'`, "", errexit, cwd)
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s %s set-e=%v: status=%d stdout=%q stderr=%q want=%q", shell, label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	missing := filepath.Join(root, "missing")
	for _, shell := range shells {
		check(shell, "file1 newer", newer, older, "", true)
		check(shell, "file1 older", older, newer, "", false)
		check(shell, "equal timestamps", equalA, equalB, "", false)
		check(shell, "same file", newer, newer, "", false)
		check(shell, "newer directory", newerDir, olderDir, "", true)
		check(shell, "older directory", olderDir, newerDir, "", false)
		check(shell, "symlink to newer target", newerLink, older, "", true)
		check(shell, "symlink to older target", olderLink, newer, "", false)
		check(shell, "newer target versus symlink to older", newer, olderLink, "", true)
		check(shell, "relative paths", "newer", "older", root, true)
		for _, name := range special {
			check(shell, "special "+name, filepath.Join(root, name), older, "", true)
		}
		check(shell, "relative leading hyphen", "-leading", "older", root, true)
		check(shell, "relative literal hyphen", "-", "older", root, true)
		for _, tc := range []struct{ label, first, second string }{
			{"right missing", newer, missing},
			{"left missing", missing, older},
			{"both missing", missing, filepath.Join(root, "another-missing")},
			{"left empty", "", older},
			{"right empty", newer, ""},
			{"both empty", "", ""},
			{"broken link left", brokenLink, older},
			{"broken link right", newer, brokenLink},
		} {
			measured := probe(shell, `[ "$filePath1" -nt "$filePath2" ]`, tc.first, tc.second, "") == 0
			t.Logf("%s %s: -nt=%v", shell, tc.label, measured)
			check(shell, tc.label, tc.first, tc.second, "", measured)
		}
		if got := probe(shell, `[ "$filePath1" -ot "$filePath2" ]`, older, newer, ""); got != 0 {
			t.Fatalf("%s -ot did not detect older file: status=%d", shell, got)
		}
		if got := probe(shell, `[ "$filePath1" -ot "$filePath2" ]`, newer, older, ""); got != 1 {
			t.Fatalf("%s -ot reversed ordering: status=%d", shell, got)
		}
		t.Logf("%s -ot with missing left operand: %v", shell, probe(shell, `[ "$filePath1" -ot "$filePath2" ]`, missing, newer, "") == 0)
		t.Logf("%s -ot with missing right operand: %v", shell, probe(shell, `[ "$filePath1" -ot "$filePath2" ]`, older, missing, "") == 0)
		if got := probe(shell, `[ -e "$filePath1" ]`, newer, "", ""); got != 0 {
			t.Fatalf("%s -e did not find existing file: status=%d", shell, got)
		}
		if got := probe(shell, `[ -e "$filePath1" ]`, missing, "", ""); got != 1 {
			t.Fatalf("%s -e found missing file: status=%d", shell, got)
		}
		if got := probe(shell, `[ -e "$filePath2" ]`, "", older, ""); got != 0 {
			t.Fatalf("%s -e did not find second existing file: status=%d", shell, got)
		}
		if got := probe(shell, `[ -e "$filePath2" ]`, "", missing, ""); got != 1 {
			t.Fatalf("%s -e found second missing file: status=%d", shell, got)
		}
		out, status, stderr := run(shell, newer, older, defaultAction, "", false, "")
		wantDefault := "Path \"" + newer + "\" is newer than path \"" + older + "\"\nSTATUS:0\n"
		if status != 0 || out != wantDefault || stderr != "" {
			t.Fatalf("%s default action: status=%d stdout=%q stderr=%q want=%q", shell, status, out, stderr, wantDefault)
		}
		for _, errexit := range []bool{false, true} {
			out, status, stderr = run(shell, newer, older, "false", "", errexit, "")
			if errexit {
				if status != 1 || out != "" || stderr != "" {
					t.Fatalf("%s set -e action failure: status=%d stdout=%q stderr=%q", shell, status, out, stderr)
				}
			} else if status != 0 || out != "STATUS:1\n" || stderr != "" {
				t.Fatalf("%s action failure: status=%d stdout=%q stderr=%q", shell, status, out, stderr)
			}
			out, status, stderr = run(shell, older, newer, "false", "", errexit, "")
			if status != 0 || out != "STATUS:0\n" || stderr != "" {
				t.Fatalf("%s false branch: status=%d stdout=%q stderr=%q", shell, status, out, stderr)
			}
		}
		out, status, stderr = run(shell, newer, older, `marker=after; printf '%s\n' "$marker"`, `marker=before`, false, "")
		if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
			t.Fatalf("%s caller shell scope: status=%d stdout=%q stderr=%q", shell, status, out, stderr)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
}

func TestIfFile1OlderThanFile2_65HistoricalContractAndBehavior(t *testing.T) {
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-file1-older-than-file2"]
	historicalBody := []any{
		`if [ "${1|/path/to/file1,${filePath1}|}" -ot "${2|/path/to/file2,${filePath2}|}" ]; then`,
		"\t${3:echo \"Path \\\"${1}\\\" is older than path \\\"${2}\\\"\"}",
		"fi\n",
	}
	if item.Prefix != "if file older" || item.Description != "check if file1 is older than file2" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-65 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-file1-older-than-file2.json")
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "unchanged historical entry-65 source", source, "f99a982654620f5a9ae2153c1bd8b7de8b943be02b9b084b398c8e3f2bd625c5")
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source differs from generated object: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	first := `${1|/path/to/file1,${filePath1}|}`
	second := `${2|/path/to/file2,${filePath2}|}`
	action := `${3:echo "Path \"${1}\" is older than path \"${2}\""}`
	for _, placeholder := range []string{first, second, action} {
		if strings.Count(body, placeholder) != 1 {
			t.Fatalf("placeholder %q absent or repeated in %q", placeholder, body)
		}
	}
	expand := func(path1, path2, branch string) string {
		return strings.Replace(strings.Replace(strings.Replace(body, first, path1, 1), second, path2, 1), action, branch, 1)
	}
	defaultAction := `echo "Path \"${filePath1}\" is older than path \"${filePath2}\""`
	for _, tc := range []struct{ name, path1, path2, branch string }{
		{"variable", `${filePath1}`, `${filePath2}`, defaultAction},
		{"literal", "/path/to/file1", "/path/to/file2", `echo "Path \"/path/to/file1\" is older than path \"/path/to/file2\""`},
	} {
		want := "if [ \"" + tc.path1 + "\" -ot \"" + tc.path2 + "\" ]; then\n\t" + tc.branch + "\nfi\n"
		if got := expand(tc.path1, tc.path2, tc.branch); got != want {
			t.Fatalf("%s expansion = %q, want %q", tc.name, got, want)
		}
	}
	root := t.TempDir()
	base := time.Unix(1600000000, 0)
	makeFile := func(name string, at time.Time) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, nil, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, at, at); err != nil {
			t.Fatal(err)
		}
		return path
	}
	old := makeFile("old", base)
	newer := makeFile("new", base.Add(10*time.Second))
	equal := makeFile("equal", base)
	for _, tc := range []struct {
		path string
		want time.Time
	}{{old, base}, {newer, base.Add(10 * time.Second)}, {equal, base}} {
		info, err := os.Stat(tc.path)
		if err != nil || !info.ModTime().Equal(tc.want) {
			t.Fatalf("fixture mtime %q: info=%v err=%v", tc.path, info, err)
		}
	}
	oldDir, newDir := filepath.Join(root, "old-dir"), filepath.Join(root, "new-dir")
	for _, tc := range []struct {
		path string
		at   time.Time
	}{{oldDir, base}, {newDir, base.Add(10 * time.Second)}} {
		if err := os.Mkdir(tc.path, 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(tc.path, tc.at, tc.at); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(tc.path)
		if err != nil || !info.ModTime().Equal(tc.at) {
			t.Fatalf("directory mtime %q: info=%v err=%v", tc.path, info, err)
		}
	}
	oldLink, newLink, broken := filepath.Join(root, "old-link"), filepath.Join(root, "new-link"), filepath.Join(root, "broken")
	for _, tc := range []struct{ link, target string }{{oldLink, old}, {newLink, newer}, {broken, filepath.Join(root, "no-target")}} {
		if err := os.Symlink(tc.target, tc.link); err != nil {
			t.Fatal(err)
		}
	}
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		makeFile(name, base)
	}
	missing := filepath.Join(root, "missing")
	shells := []string{"bash", "sh"}
	if _, err := exec.LookPath("zsh"); err == nil {
		shells = append(shells, "zsh")
	}
	run := func(shell, code, path1, path2, cwd string) (string, int, string) {
		t.Helper()
		cmd := exec.Command(shell, "-c", code)
		cmd.Env = append(os.Environ(), "filePath1="+path1, "filePath2="+path2)
		if cwd != "" {
			cmd.Dir = cwd
		}
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		return string(out), exitCode(err), stderr.String()
	}
	probe := func(shell, expr, path1, path2, cwd string) bool {
		t.Helper()
		_, status, stderr := run(shell, expr, path1, path2, cwd)
		if status > 1 || stderr != "" {
			t.Fatalf("%s probe %q: status=%d stderr=%q", shell, expr, status, stderr)
		}
		return status == 0
	}
	check := func(shell, label, path1, path2, cwd string, want bool) {
		t.Helper()
		if got := probe(shell, `[ "$filePath1" -ot "$filePath2" ]`, path1, path2, cwd); got != want {
			t.Fatalf("%s %s direct predicate=%v, want %v", shell, label, got, want)
		}
		for _, errexit := range []bool{false, true} {
			code := expand(`${filePath1}`, `${filePath2}`, `printf 'TRUE\n'`)
			if errexit {
				code = "set -e\n" + code
			}
			out, status, stderr := run(shell, code+`printf 'STATUS:%s\n' "$?"`+"\n", path1, path2, cwd)
			wantOut := "STATUS:0\n"
			if want {
				wantOut = "TRUE\nSTATUS:0\n"
			}
			if status != 0 || out != wantOut || stderr != "" {
				t.Fatalf("%s %s set-e=%v: status=%d out=%q stderr=%q want=%q", shell, label, errexit, status, out, stderr, wantOut)
			}
		}
	}
	for _, shell := range shells {
		for _, tc := range []struct {
			label, left, right, cwd string
			want                    bool
		}{
			{"older file", old, newer, "", true}, {"newer file", newer, old, "", false},
			{"equal", old, equal, "", false}, {"same file", old, old, "", false},
			{"older directory", oldDir, newDir, "", true}, {"newer directory", newDir, oldDir, "", false},
			{"old symlink", oldLink, newer, "", true}, {"new symlink", newLink, old, "", false},
			{"symlink on right", old, newLink, "", true},
			{"relative", "old", "new", root, true},
			{"relative leading hyphen", "-leading", "new", root, true},
			{"relative literal hyphen", "-", "new", root, true},
		} {
			check(shell, tc.label, tc.left, tc.right, tc.cwd, tc.want)
		}
		for _, name := range special {
			check(shell, "special "+name, filepath.Join(root, name), newer, "", true)
		}
		for _, tc := range []struct{ label, left, right string }{
			{"missing left", missing, newer}, {"missing right", old, missing},
			{"both missing", missing, filepath.Join(root, "other-missing")},
			{"empty left", "", newer}, {"empty right", old, ""}, {"both empty", "", ""},
			{"broken left", broken, newer}, {"broken right", old, broken},
		} {
			measured := probe(shell, `[ "$filePath1" -ot "$filePath2" ]`, tc.left, tc.right, "")
			t.Logf("%s %s: -ot=%v", shell, tc.label, measured)
			if shell == "bash" && tc.label == "missing left" && !measured {
				t.Fatal("Bash must treat a missing left operand as older than an existing right operand")
			}
			check(shell, tc.label, tc.left, tc.right, "", measured)
		}
		for _, tc := range []struct {
			expr, left, right string
			want              bool
		}{
			{`[ "$filePath1" -nt "$filePath2" ]`, newer, old, true},
			{`[ "$filePath1" -nt "$filePath2" ]`, old, newer, false},
			{`[ -e "$filePath1" ]`, old, "", true},
			{`[ -e "$filePath1" ]`, missing, "", false},
			{`[ -e "$filePath2" ]`, "", newer, true},
			{`[ -e "$filePath2" ]`, "", missing, false},
		} {
			if got := probe(shell, tc.expr, tc.left, tc.right, ""); got != tc.want {
				t.Fatalf("%s probe %q=%v, want %v", shell, tc.expr, got, tc.want)
			}
		}
		t.Logf("%s missing right: -nt=%v", shell, probe(shell, `[ "$filePath1" -nt "$filePath2" ]`, newer, missing, ""))
		out, status, stderr := run(shell, expand(`${filePath1}`, `${filePath2}`, defaultAction), old, newer, "")
		want := "Path \"" + old + "\" is older than path \"" + newer + "\"\n"
		if status != 0 || out != want || stderr != "" {
			t.Fatalf("%s default action: status=%d out=%q stderr=%q want=%q", shell, status, out, stderr, want)
		}
		for _, errexit := range []bool{false, true} {
			code := expand(`${filePath1}`, `${filePath2}`, "false")
			if errexit {
				code = "set -e\n" + code
			}
			out, status, stderr = run(shell, code+`printf 'STATUS:%s\n' "$?"`+"\n", old, newer, "")
			if errexit && (status != 1 || out != "" || stderr != "") {
				t.Fatalf("%s set-e true branch failure: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
			if !errexit && (status != 0 || out != "STATUS:1\n" || stderr != "") {
				t.Fatalf("%s true branch failure: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
			out, status, stderr = run(shell, code+`printf 'STATUS:%s\n' "$?"`+"\n", newer, old, "")
			if status != 0 || out != "STATUS:0\n" || stderr != "" {
				t.Fatalf("%s false branch: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
		}
		out, status, stderr = run(shell, "marker=before\n"+expand(`${filePath1}`, `${filePath2}`, `marker=after; printf '%s\n' "$marker"`)+`printf 'STATUS:%s\n' "$?"`+"\n", old, newer, "")
		if status != 0 || out != "after\nSTATUS:0\n" || stderr != "" {
			t.Fatalf("%s caller scope: status=%d out=%q stderr=%q", shell, status, out, stderr)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
}

func TestIfFilesAreEqual66GeneratedBehavior(t *testing.T) {
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "iteration-68 artifact while checking entry 66", data, currentV7Iteration68SHA256)
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-files-are-equal"]
	if item.Prefix != "if file =" || item.Description != "check if files are equal" {
		t.Fatalf("entry-66 metadata: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-files-are-equal.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("entry-66 source unexpectedly ends in a newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("source and generated entry differ: %v", err)
	}
	lines := anyStrings(t, item.Body.([]any))
	if len(lines) != 21 || lines[0] != "if (" || lines[10] != "); then" || lines[12] != "else" || lines[20] != "fi\n" {
		t.Fatalf("entry-66 body structure: %#v", lines)
	}
	body := strings.Join(lines, "\n")
	first := `${1|/path/to/file1,${filePath1}|}`
	second := `${2|/path/to/file2,${filePath2}|}`
	action := `${3:echo "File \"${1}\" is equal to file \"${2}\""}`
	for _, token := range []string{first, second, action} {
		if strings.Count(body, token) != 1 {
			t.Fatalf("placeholder %q absent or repeated in %q", token, body)
		}
	}
	expand := func(path1, path2, branch string) string {
		return strings.Replace(strings.Replace(strings.Replace(body, first, path1, 1), second, path2, 1), action, branch, 1)
	}
	variableAction := `echo "File \"${filePath1}\" is equal to file \"${filePath2}\""`
	variableScript := expand(`${filePath1}`, `${filePath2}`, variableAction)
	if !strings.Contains(variableScript, `file1="${filePath1}"`) || !strings.Contains(variableScript, `file2="${filePath2}"`) || !strings.Contains(variableScript, "\t"+variableAction+"\n") {
		t.Fatalf("variable expansion = %q", variableScript)
	}
	literalAction := `echo "File \"/path/to/file1\" is equal to file \"/path/to/file2\""`
	literalScript := expand("/path/to/file1", "/path/to/file2", literalAction)
	if !strings.Contains(literalScript, `file1="/path/to/file1"`) || !strings.Contains(literalScript, `file2="/path/to/file2"`) || !strings.Contains(literalScript, "\t"+literalAction+"\n") {
		t.Fatalf("literal expansion = %q", literalScript)
	}
	root := t.TempDir()
	base := time.Unix(1600000000, 0)
	makeFile := func(name string, content []byte, at time.Time) string {
		t.Helper()
		path := filepath.Join(root, name)
		if err := os.WriteFile(path, content, 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, at, at); err != nil {
			t.Fatal(err)
		}
		return path
	}
	content := []byte("same content\n")
	a := makeFile("a", content, base)
	b := makeFile("b", content, base.Add(10*time.Second))
	c := makeFile("c", []byte("other content"), base)
	if len(content) != len([]byte("other content")) {
		t.Fatal("equal-size different-content fixture has wrong size")
	}
	for _, tc := range []struct {
		path string
		at   time.Time
	}{{a, base}, {b, base.Add(10 * time.Second)}, {c, base}} {
		info, err := os.Stat(tc.path)
		if err != nil || !info.ModTime().Equal(tc.at) {
			t.Fatalf("fixture mtime %q: info=%v err=%v", tc.path, info, err)
		}
	}
	emptyA := makeFile("empty-a", nil, base)
	emptyB := makeFile("empty-b", nil, base)
	binaryA := makeFile("binary-a", []byte{'a', 0, 'b'}, base)
	binaryB := makeFile("binary-b", []byte{'a', 0, 'b'}, base.Add(10*time.Second))
	binaryC := makeFile("binary-c", []byte{'a', 0, 'c'}, base)
	noFinalNewlineA := makeFile("no-final-newline-a", []byte("line"), base)
	noFinalNewlineB := makeFile("no-final-newline-b", []byte("line"), base)
	noFinalNewlineC := makeFile("no-final-newline-c", []byte("line\n"), base)
	hardLink := filepath.Join(root, "hard-link")
	if err := os.Link(a, hardLink); err != nil {
		t.Fatal(err)
	}
	linkA, linkA2, linkB, broken := filepath.Join(root, "link-a"), filepath.Join(root, "link-a2"), filepath.Join(root, "link-b"), filepath.Join(root, "broken")
	for _, tc := range []struct{ link, target string }{{linkA, a}, {linkA2, a}, {linkB, b}, {broken, filepath.Join(root, "no-target")}} {
		if err := os.Symlink(tc.target, tc.link); err != nil {
			t.Fatal(err)
		}
	}
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		makeFile(name, content, base)
	}
	missing := filepath.Join(root, "missing")
	probeCmp := func(args ...string) (int, string, string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "cmp", append([]string{"-s"}, args...)...)
		cmd.Dir = root
		cmd.Stdin = strings.NewReader("different stdin")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if ctx.Err() != nil {
			t.Fatal("direct cmp probe timed out")
		}
		return exitCode(err), string(out), stderr.String()
	}
	for _, tc := range []struct {
		label, first, second string
		want                 int
	}{
		{"identical", "a", "b", 0}, {"different", "a", "c", 1},
		{"missing", "missing", "a", 2}, {"historical leading hyphen", "-leading", "a", 2},
		{"historical dash means stdin", "-", "a", 1},
	} {
		status, out, stderr := probeCmp(tc.first, tc.second)
		if status != tc.want || out != "" {
			t.Fatalf("direct cmp %s: status=%d out=%q stderr=%q", tc.label, status, out, stderr)
		}
		if tc.label == "historical leading hyphen" && stderr == "" {
			t.Fatal("historical leading hyphen did not produce an option diagnostic")
		}
		if tc.label == "missing" && stderr != "" {
			t.Fatalf("silent missing-file cmp emitted %q", stderr)
		}
	}
	directory := filepath.Join(root, "directory")
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	unreadable := makeFile("unreadable", content, base)
	if err := os.Chmod(unreadable, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(unreadable, 0o600) })
	where := func(name string) string { return filepath.Join(root, name) }
	stdoutStatus := `printf 'STATUS:%s\n' "$?"` + "\n"
	run := func(shell, script, path1, path2, cwd string, stdin string) (string, int, string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, shell, "-c", script)
		cmd.Env = append(os.Environ(), "filePath1="+path1, "filePath2="+path2)
		cmd.Stdin = strings.NewReader(stdin)
		if cwd != "" {
			cmd.Dir = cwd
		}
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if ctx.Err() != nil {
			t.Fatalf("%s timed out running comparison", shell)
		}
		return string(out), exitCode(err), stderr.String()
	}
	check := func(shell, label, path1, path2, cwd string, equal bool) {
		t.Helper()
		want := "STATUS:0\n"
		if equal {
			want = "TRUE\nSTATUS:0\n"
		}
		for _, errexit := range []bool{false, true} {
			script := expand(`${filePath1}`, `${filePath2}`, `printf 'TRUE\n'`)
			if errexit {
				script = "set -e\n" + script
			}
			out, status, stderr := run(shell, script+stdoutStatus, path1, path2, cwd, "different stdin")
			if status != 0 || out != want || stderr != "" {
				t.Fatalf("%s %s set-e=%v: status=%d out=%q stderr=%q want=%q", shell, label, errexit, status, out, stderr, want)
			}
		}
	}
	shells := []string{"bash", "sh"}
	if _, err := exec.LookPath("zsh"); err == nil {
		shells = append(shells, "zsh")
	}
	for _, shell := range shells {
		for _, tc := range []struct {
			label, left, right, cwd string
			equal                   bool
		}{
			{"identical content", a, b, "", true}, {"different content", a, c, "", false},
			{"equal size, different content", b, c, "", false},
			{"equal content, different mtime", a, b, "", true},
			{"different content, equal mtime", a, c, "", false},
			{"empty files", emptyA, emptyB, "", true},
			{"binary equal", binaryA, binaryB, "", true}, {"binary different", binaryA, binaryC, "", false},
			{"no final newline equal", noFinalNewlineA, noFinalNewlineB, "", true},
			{"no final newline different", noFinalNewlineA, noFinalNewlineC, "", false},
			{"same path", a, a, "", true}, {"hard link", a, hardLink, "", true},
			{"symlinks same target", linkA, linkA2, "", true},
			{"symlinks different targets, same contents", linkA, linkB, "", true},
			{"relative paths", "a", "b", root, true},
			{"relative leading hyphen", "-leading", "a", root, true},
			{"relative leading hyphen right", "a", "-leading", root, true},
			{"relative literal hyphen", "-", "a", root, true},
			{"relative literal hyphen right", "a", "-", root, true},
		} {
			check(shell, tc.label, tc.left, tc.right, tc.cwd, tc.equal)
		}
		for _, name := range special {
			check(shell, "special "+name, where(name), a, "", true)
		}
		for _, tc := range []struct{ label, left, right string }{
			{"missing left", missing, a}, {"missing right", a, missing},
			{"both missing", missing, where("other-missing")},
			{"empty left", "", a}, {"empty right", a, ""}, {"both empty", "", ""},
			{"broken symlink left", broken, a}, {"broken symlink right", a, broken},
			{"directory left", directory, a}, {"directory right", a, directory},
		} {
			for _, errexit := range []bool{false, true} {
				script := expand(`${filePath1}`, `${filePath2}`, `printf 'TRUE\n'`)
				if errexit {
					script = "set -e\n" + script
				}
				out, status, stderr := run(shell, script+stdoutStatus, tc.left, tc.right, "", "")
				wantStderr := "cmp: comparison failed\n"
				if strings.HasPrefix(tc.label, "directory") {
					if !strings.Contains(stderr, "Is a directory\n") || !strings.HasSuffix(stderr, wantStderr) {
						t.Fatalf("%s %s directory diagnostic = %q", shell, tc.label, stderr)
					}
					wantStderr = stderr
				}
				if errexit {
					if status != 2 || out != "" || stderr != wantStderr {
						t.Fatalf("%s %s set-e error: status=%d out=%q stderr=%q", shell, tc.label, status, out, stderr)
					}
				} else if status != 0 || out != "STATUS:2\n" || stderr != wantStderr {
					t.Fatalf("%s %s error: status=%d out=%q stderr=%q", shell, tc.label, status, out, stderr)
				}
			}
		}
		if os.Geteuid() != 0 {
			out, status, stderr := run(shell, expand(`${filePath1}`, `${filePath2}`, `printf 'TRUE\n'`)+stdoutStatus, unreadable, a, "", "")
			if status != 0 || out != "STATUS:2\n" || stderr != "cmp: comparison failed\n" {
				t.Fatalf("%s unreadable file: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
		}
		out, status, stderr := run(shell, expand(`${filePath1}`, `${filePath2}`, variableAction), a, b, "", "")
		want := "File \"" + a + "\" is equal to file \"" + b + "\"\n"
		if status != 0 || out != want || stderr != "" {
			t.Fatalf("%s default action: status=%d out=%q stderr=%q want=%q", shell, status, out, stderr, want)
		}
		for _, errexit := range []bool{false, true} {
			script := expand(`${filePath1}`, `${filePath2}`, "false")
			if errexit {
				script = "set -e\n" + script
			}
			out, status, stderr = run(shell, script+stdoutStatus, a, b, "", "")
			if errexit && (status != 1 || out != "" || stderr != "") {
				t.Fatalf("%s set-e action failure: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
			if !errexit && (status != 0 || out != "STATUS:1\n" || stderr != "") {
				t.Fatalf("%s action failure: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
			out, status, stderr = run(shell, script+stdoutStatus, a, c, "", "")
			if status != 0 || out != "STATUS:0\n" || stderr != "" {
				t.Fatalf("%s false branch: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
		}
		out, status, stderr = run(shell, "marker=before\n"+expand(`${filePath1}`, `${filePath2}`, `marker=after; printf '%s\n' "$marker"`)+`printf 'SCOPE:%s:%s:%s\n' "$marker" "${file1-unset}" "${file2-unset}"`+"\n", a, b, "", "")
		if status != 0 || out != "after\nSCOPE:after:unset:unset\n" || stderr != "" {
			t.Fatalf("%s caller scope: status=%d out=%q stderr=%q", shell, status, out, stderr)
		}
	}
	if _, err := os.Stat(where("injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
}

func TestIfPathExists67HistoricalContractAndBehavior(t *testing.T) {
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.if-path-exists"]
	historicalBody := []any{
		`if [ -e "${1|/path/to/something,${pathToSomething}|}" ]; then`,
		"\t${2:echo \"Path \\\"${1}\\\" exists\"}",
		"fi\n",
	}
	if item.Prefix != "if path exists" || item.Description != "if path exists (file, directory, link...)" || !reflect.DeepEqual(item.Body, historicalBody) {
		t.Fatalf("entry-67 historical contract: %#v", item)
	}
	source, err := os.ReadFile("nsroot/filesystem/if-path-exists.json")
	if err != nil {
		t.Fatal(err)
	}
	assertSHA256(t, "unchanged historical entry-67 source", source, "189d11b2dd683dec681c143bf8ffc6c685931092cfa45050db28dcdf432ce232")
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("historical source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("historical source and generated object differ: %v", err)
	}
	body := strings.Join(anyStrings(t, item.Body.([]any)), "\n")
	pathPlaceholder := `${1|/path/to/something,${pathToSomething}|}`
	actionPlaceholder := `${2:echo "Path \"${1}\" exists"}`
	if strings.Count(body, pathPlaceholder) != 1 || strings.Count(body, actionPlaceholder) != 1 {
		t.Fatalf("placeholder spelling or order changed: %q", body)
	}
	expand := func(path, action string) string {
		return strings.Replace(strings.Replace(body, pathPlaceholder, path, 1), actionPlaceholder, action, 1)
	}
	defaultAction := `echo "Path \"${pathToSomething}\" exists"`
	if got, want := expand(`${pathToSomething}`, defaultAction), "if [ -e \"${pathToSomething}\" ]; then\n\t"+defaultAction+"\nfi\n"; got != want {
		t.Fatalf("variable-choice expansion = %q, want %q", got, want)
	}
	literalAction := `echo "Path \"/path/to/something\" exists"`
	if got, want := expand("/path/to/something", literalAction), "if [ -e \"/path/to/something\" ]; then\n\t"+literalAction+"\nfi\n"; got != want {
		t.Fatalf("literal-choice expansion = %q, want %q", got, want)
	}
	root := t.TempDir()
	regular := filepath.Join(root, "regular")
	if err := os.WriteFile(regular, []byte("data"), 0o600); err != nil {
		t.Fatal(err)
	}
	directory := filepath.Join(root, "directory")
	if err := os.Mkdir(directory, 0o700); err != nil {
		t.Fatal(err)
	}
	linkFile, linkDirectory, broken, chain := filepath.Join(root, "link-file"), filepath.Join(root, "link-directory"), filepath.Join(root, "broken"), filepath.Join(root, "chain")
	for _, tc := range []struct{ link, target string }{
		{linkFile, regular}, {linkDirectory, directory}, {broken, filepath.Join(root, "absent-target")}, {chain, linkFile},
	} {
		if err := os.Symlink(tc.target, tc.link); err != nil {
			t.Fatal(err)
		}
	}
	intermediate := filepath.Join(root, "intermediate")
	if err := os.Symlink(filepath.Join(root, "absent-directory"), intermediate); err != nil {
		t.Fatal(err)
	}
	danglingIntermediate := filepath.Join(intermediate, "child")
	missing := filepath.Join(root, "missing")
	special := []string{"space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"}
	for _, name := range special {
		if err := os.WriteFile(filepath.Join(root, name), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	lockedDirectory := filepath.Join(root, "locked")
	if err := os.Mkdir(lockedDirectory, 0o700); err != nil {
		t.Fatal(err)
	}
	lockedChild := filepath.Join(lockedDirectory, "child")
	if err := os.WriteFile(lockedChild, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(lockedDirectory, 0); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chmod(lockedDirectory, 0o700) })
	_, lockedErr := os.Stat(lockedChild)
	t.Logf("effective uid=%d, stat through nonsearchable parent: %v", os.Geteuid(), lockedErr)
	shells := []string{"sh", "bash"}
	if _, err := exec.LookPath("zsh"); err == nil {
		shells = append(shells, "zsh")
	}
	run := func(shell, script, path, cwd string) (string, int, string) {
		t.Helper()
		ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, shell, "-c", script)
		cmd.Env = append(os.Environ(), "pathToSomething="+path)
		if cwd != "" {
			cmd.Dir = cwd
		}
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if ctx.Err() != nil {
			t.Fatalf("%s path test timed out", shell)
		}
		return string(out), exitCode(err), stderr.String()
	}
	probe := func(shell, primary, path, cwd string) bool {
		t.Helper()
		_, status, stderr := run(shell, `[ `+primary+` "$pathToSomething" ]`, path, cwd)
		if status > 1 || stderr != "" {
			t.Fatalf("%s %s probe %q: status=%d stderr=%q", shell, primary, path, status, stderr)
		}
		return status == 0
	}
	check := func(shell, label, path, cwd string, exists bool) {
		t.Helper()
		if got := probe(shell, "-e", path, cwd); got != exists {
			t.Fatalf("%s %s direct -e=%v, want %v", shell, label, got, exists)
		}
		for _, errexit := range []bool{false, true} {
			script := expand(`${pathToSomething}`, `printf 'TRUE\n'`)
			if errexit {
				script = "set -e\n" + script
			}
			out, status, stderr := run(shell, script+`printf 'STATUS:%s\n' "$?"`+"\n", path, cwd)
			want := "STATUS:0\n"
			if exists {
				want = "TRUE\nSTATUS:0\n"
			}
			if status != 0 || out != want || stderr != "" {
				t.Fatalf("%s %s set-e=%v: status=%d out=%q stderr=%q want=%q", shell, label, errexit, status, out, stderr, want)
			}
		}
	}
	for _, shell := range shells {
		for _, tc := range []struct {
			label, path, cwd string
			exists           bool
		}{
			{"regular", regular, "", true}, {"directory", directory, "", true},
			{"symlink to file", linkFile, "", true}, {"symlink to directory", linkDirectory, "", true},
			{"broken symlink", broken, "", false}, {"symlink chain", chain, "", true},
			{"dangling intermediate", danglingIntermediate, "", false},
			{"missing", missing, "", false}, {"empty", "", "", false},
			{"relative regular", "regular", root, true},
			{"relative leading hyphen", "-leading", root, true},
			{"relative literal hyphen", "-", root, true},
			{"directory trailing slash", directory + "/", "", true},
			{"regular trailing slash", regular + "/", "", false},
			{"symlink directory trailing slash", linkDirectory + "/", "", true},
			{"symlink file trailing slash", linkFile + "/", "", false},
			{"broken symlink trailing slash", broken + "/", "", false},
		} {
			check(shell, tc.label, tc.path, tc.cwd, tc.exists)
		}
		for _, name := range special {
			check(shell, "special "+name, filepath.Join(root, name), "", true)
		}
		for _, tc := range []struct {
			label, path   string
			e, L, h, f, d bool
		}{
			{"regular", regular, true, false, false, true, false},
			{"directory", directory, true, false, false, false, true},
			{"symlink file", linkFile, true, true, true, true, false},
			{"symlink directory", linkDirectory, true, true, true, false, true},
			{"broken symlink", broken, false, true, true, false, false},
			{"symlink chain", chain, true, true, true, true, false},
			{"dangling intermediate", danglingIntermediate, false, false, false, false, false},
			{"missing", missing, false, false, false, false, false},
			{"empty", "", false, false, false, false, false},
		} {
			for _, primary := range []struct {
				op   string
				want bool
			}{{"-e", tc.e}, {"-L", tc.L}, {"-h", tc.h}, {"-f", tc.f}, {"-d", tc.d}} {
				if got := probe(shell, primary.op, tc.path, ""); got != primary.want {
					t.Fatalf("%s %s %s=%v, want %v", shell, tc.label, primary.op, got, primary.want)
				}
			}
		}
		lockedResolved := lockedErr == nil
		check(shell, "nonsearchable parent", lockedChild, "", lockedResolved)
		t.Logf("%s nonsearchable parent: -e=%v", shell, probe(shell, "-e", lockedChild, ""))
		out, status, stderr := run(shell, expand(`${pathToSomething}`, defaultAction), regular, "")
		want := "Path \"" + regular + "\" exists\n"
		if status != 0 || out != want || stderr != "" {
			t.Fatalf("%s default action: status=%d out=%q stderr=%q want=%q", shell, status, out, stderr, want)
		}
		for _, errexit := range []bool{false, true} {
			script := expand(`${pathToSomething}`, "false")
			if errexit {
				script = "set -e\n" + script
			}
			out, status, stderr = run(shell, script+`printf 'STATUS:%s\n' "$?"`+"\n", regular, "")
			if errexit && (status != 1 || out != "" || stderr != "") {
				t.Fatalf("%s set-e action failure: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
			if !errexit && (status != 0 || out != "STATUS:1\n" || stderr != "") {
				t.Fatalf("%s action failure: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
			out, status, stderr = run(shell, script+`printf 'STATUS:%s\n' "$?"`+"\n", missing, "")
			if status != 0 || out != "STATUS:0\n" || stderr != "" {
				t.Fatalf("%s false branch: status=%d out=%q stderr=%q", shell, status, out, stderr)
			}
		}
		out, status, stderr = run(shell, "marker=before\n"+expand(`${pathToSomething}`, `marker=after; printf '%s\n' "$marker"`)+`printf 'STATUS:%s\n' "$marker"`+"\n", regular, "")
		if status != 0 || out != "after\nSTATUS:after\n" || stderr != "" {
			t.Fatalf("%s caller scope: status=%d out=%q stderr=%q", shell, status, out, stderr)
		}
	}
	if _, err := os.Stat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("injection marker appeared: %v", err)
	}
}

func TestRemoveFilesOlderNewer68GeneratedSafeBehavior(t *testing.T) {
	data, err := os.ReadFile(snippetOutputPath)
	if err != nil {
		t.Fatal(err)
	}
	var snippets map[string]Snippet
	if err := json.Unmarshal(data, &snippets); err != nil {
		t.Fatal(err)
	}
	item := snippets["filesystem.remove-files-older-newer-than"]
	if item.Prefix != "remove old/new files/directories" || item.Description != "find and remove files(f)/directories(d) older(+)/newer(-) than x days" {
		t.Fatalf("entry-68 metadata = %#v", item)
	}
	body, ok := item.Body.(string)
	if !ok {
		t.Fatalf("entry-68 body representation = %T, want string", item.Body)
	}
	source, err := os.ReadFile("nsroot/filesystem/remove-files-older-newer-than.json")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.HasSuffix(source, []byte("\n")) {
		t.Fatal("entry-68 source unexpectedly ends in newline")
	}
	var sourceItem Snippet
	if err := json.Unmarshal(source, &sourceItem); err != nil || !reflect.DeepEqual(sourceItem, item) {
		t.Fatalf("entry-68 source and generated object differ: %v", err)
	}
	pathToken := `${1|/path/to/directory/,${pathToDirectory}|}`
	typeToken := `${2|f,d|}`
	signToken := `${3|-,+|}`
	daysToken := `${4:days}`
	for _, token := range []string{pathToken, typeToken, signToken, daysToken} {
		if strings.Count(body, token) != 1 {
			t.Fatalf("placeholder %q absent or repeated in %q", token, body)
		}
	}
	if !(strings.Index(body, pathToken) < strings.Index(body, typeToken) && strings.Index(body, typeToken) < strings.Index(body, signToken) && strings.Index(body, signToken) < strings.Index(body, daysToken)) {
		t.Fatal("entry-68 placeholder order changed")
	}
	expand := func(path, kind, sign, days string) string {
		return strings.Replace(strings.Replace(strings.Replace(strings.Replace(body, pathToken, path, 1), typeToken, kind, 1), signToken, sign, 1), daysToken, days, 1)
	}
	for _, tc := range []struct{ path, kind, sign, days, want string }{
		{`${pathToDirectory}`, "f", "-", "1", `search_dir="${pathToDirectory}"`},
		{"/path/to/directory/", "d", "+", "2", `search_dir="/path/to/directory/"`},
	} {
		got := expand(tc.path, tc.kind, tc.sign, tc.days)
		if !strings.Contains(got, tc.want) || !strings.Contains(got, "CDPATH='' cd -P \"$search_dir\" || exit 2") || !strings.Contains(got, "for entry in ./*; do") || !strings.Contains(got, "find \"$entry\" -type "+tc.kind+" -mtime "+tc.sign+tc.days+" -delete") || !strings.HasSuffix(got, ")\n") {
			t.Fatalf("generated expansion = %q", got)
		}
	}
	root, err := os.MkdirTemp("/tmp", "shellman-iteration-68-")
	if err != nil {
		t.Fatal(err)
	}
	rootCanonical, err := filepath.EvalSymlinks(root)
	if err != nil {
		t.Fatal(err)
	}
	tmpCanonical, err := filepath.EvalSymlinks("/tmp")
	if err != nil {
		t.Fatal(err)
	}
	repoCanonical, err := filepath.EvalSymlinks(".")
	if err != nil {
		t.Fatal(err)
	}
	contained := func(parent, child string, proper bool) bool {
		rel, err := filepath.Rel(parent, child)
		if err != nil || filepath.IsAbs(rel) || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return false
		}
		return !proper || rel != "."
	}
	if rootCanonical == "/" || rootCanonical == tmpCanonical || rootCanonical == repoCanonical || !contained(tmpCanonical, rootCanonical, true) || contained(repoCanonical, rootCanonical, false) {
		t.Fatalf("unsafe disposable root %q (tmp=%q repo=%q)", rootCanonical, tmpCanonical, repoCanonical)
	}
	outside, err := os.CreateTemp("/tmp", "shellman-iteration-68-outside-")
	if err != nil {
		t.Fatal(err)
	}
	outsidePath := outside.Name()
	if _, err := outside.WriteString("outside sentinel"); err != nil {
		t.Fatal(err)
	}
	if err := outside.Close(); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Remove(outsidePath) })
	if contained(rootCanonical, outsidePath, false) {
		t.Fatal("outside sentinel unexpectedly lies within fixture root")
	}
	assertRoot := func() {
		t.Helper()
		info, err := os.Lstat(root)
		if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
			t.Fatalf("disposable root no longer a real directory: info=%v err=%v", info, err)
		}
		canonical, err := filepath.EvalSymlinks(root)
		if err != nil || canonical != rootCanonical || !contained(tmpCanonical, canonical, true) || canonical == repoCanonical || canonical == "/" {
			t.Fatalf("disposable root changed: %q err=%v", canonical, err)
		}
	}
	t.Cleanup(func() {
		assertRoot()
		if err := os.RemoveAll(root); err != nil {
			t.Errorf("fixture cleanup: %v", err)
		}
	})
	assertScope := func(scope, input string) {
		t.Helper()
		assertRoot()
		canonical, err := filepath.EvalSymlinks(scope)
		if err != nil || !contained(rootCanonical, canonical, true) {
			t.Fatalf("search scope outside disposable root: %q -> %q err=%v", scope, canonical, err)
		}
		info, err := os.Stat(scope)
		if err != nil || !info.IsDir() {
			t.Fatalf("search scope is not a directory: %q info=%v err=%v", scope, info, err)
		}
		inputPath := input
		if !filepath.IsAbs(inputPath) {
			inputPath = filepath.Join(root, inputPath)
		}
		inputCanonical, err := filepath.EvalSymlinks(inputPath)
		if err != nil || inputCanonical != canonical {
			t.Fatalf("user input %q escapes intended scope %q: resolved=%q err=%v", input, scope, inputCanonical, err)
		}
		err = filepath.WalkDir(scope, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if !contained(rootCanonical, filepath.Clean(path), true) {
				return fmt.Errorf("candidate outside fixture root: %q", path)
			}
			if d.Type()&os.ModeSymlink != 0 {
				target, err := os.Readlink(path)
				if err != nil {
					return err
				}
				if !filepath.IsAbs(target) {
					target = filepath.Join(filepath.Dir(path), target)
				}
				target = filepath.Clean(target)
				if !contained(rootCanonical, target, true) {
					return fmt.Errorf("fixture symlink points outside root: %q -> %q", path, target)
				}
			}
			return nil
		})
		if err != nil {
			t.Fatalf("search containment walk failed: %v", err)
		}
	}
	type inventoryEntry struct {
		mode os.FileMode
		data string
		link string
	}
	inventory := func() map[string]inventoryEntry {
		t.Helper()
		result := make(map[string]inventoryEntry)
		err := filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			rel, err := filepath.Rel(root, path)
			if err != nil {
				return err
			}
			info, err := os.Lstat(path)
			if err != nil {
				return err
			}
			entry := inventoryEntry{mode: info.Mode()}
			if info.Mode().IsRegular() {
				bytes, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				entry.data = string(bytes)
			} else if info.Mode()&os.ModeSymlink != 0 {
				entry.link, err = os.Readlink(path)
				if err != nil {
					return err
				}
			}
			result[rel] = entry
			return nil
		})
		if err != nil {
			t.Fatalf("fixture inventory failed: %v", err)
		}
		return result
	}
	writeAt := func(path, data string, at time.Time) {
		t.Helper()
		if !contained(rootCanonical, filepath.Clean(path), true) {
			t.Fatalf("fixture file outside root: %q", path)
		}
		if err := os.WriteFile(path, []byte(data), 0o600); err != nil {
			t.Fatal(err)
		}
		if err := os.Chtimes(path, at, at); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil || !info.ModTime().Equal(at) {
			t.Fatalf("fixture timestamp %q: info=%v err=%v", path, info, err)
		}
	}
	mkdir := func(path string) {
		t.Helper()
		if !contained(rootCanonical, filepath.Clean(path), true) {
			t.Fatalf("fixture directory outside root: %q", path)
		}
		if err := os.Mkdir(path, 0o700); err != nil {
			t.Fatal(err)
		}
	}
	now := time.Now().Truncate(time.Second)
	old := now.Add(-72 * time.Hour)
	newer := now.Add(-time.Hour)
	setDirTime := func(path string, at time.Time) {
		t.Helper()
		if err := os.Chtimes(path, at, at); err != nil {
			t.Fatal(err)
		}
		info, err := os.Stat(path)
		if err != nil || !info.ModTime().Equal(at) {
			t.Fatalf("directory timestamp %q: info=%v err=%v", path, info, err)
		}
	}
	run := func(label, scope, input, kind, sign, days string, errexit bool, wantStatus int, wantDeleted []string, wantDiagnostic bool) {
		t.Helper()
		assertScope(scope, input)
		for _, rel := range wantDeleted {
			if !contained(rootCanonical, filepath.Join(root, rel), true) || !contained(scope, filepath.Join(root, rel), true) {
				t.Fatalf("%s intended deletion outside search subtree: %q", label, rel)
			}
		}
		before := inventory()
		script := expand(`${pathToDirectory}`, kind, sign, days)
		if errexit {
			script = "set -e\n" + script
		}
		script += `printf 'STATUS:%s\n' "$?"` + "\n"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "sh", "-c", script)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "pathToDirectory="+input, "LC_ALL=C", "CDPATH="+filepath.Join(root, "cdpath-trap"))
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if ctx.Err() != nil {
			t.Fatalf("%s destructive test timed out", label)
		}
		status := exitCode(err)
		wantOut := fmt.Sprintf("STATUS:%d\n", wantStatus)
		wantProcessStatus := 0
		if errexit && wantStatus != 0 {
			wantOut = ""
			wantProcessStatus = wantStatus
		}
		if status != wantProcessStatus || string(out) != wantOut || (stderr.Len() != 0) != wantDiagnostic {
			t.Fatalf("%s: status=%d out=%q stderr=%q, want process=%d out=%q diagnostic=%v", label, status, out, stderr.String(), wantProcessStatus, wantOut, wantDiagnostic)
		}
		after := inventory()
		deleted := make([]string, 0)
		for rel, prior := range before {
			current, exists := after[rel]
			if !exists {
				deleted = append(deleted, rel)
			} else if !reflect.DeepEqual(prior, current) {
				t.Fatalf("%s non-target changed: %q before=%#v after=%#v", label, rel, prior, current)
			}
		}
		for rel := range after {
			if _, existed := before[rel]; !existed {
				t.Fatalf("%s unexpectedly created %q", label, rel)
			}
		}
		sort.Strings(deleted)
		sort.Strings(wantDeleted)
		if !reflect.DeepEqual(deleted, wantDeleted) {
			t.Fatalf("%s deletion set=%q, want %q", label, deleted, wantDeleted)
		}
		outsideData, err := os.ReadFile(outsidePath)
		if err != nil || string(outsideData) != "outside sentinel" {
			t.Fatalf("%s outside sentinel changed: data=%q err=%v", label, outsideData, err)
		}
	}
	cdpathTrap := filepath.Join(root, "cdpath-trap")
	mkdir(cdpathTrap)
	mkdir(filepath.Join(cdpathTrap, "target"))
	writeAt(filepath.Join(cdpathTrap, "target", "old"), "CDPATH sentinel", old)
	mainDir := filepath.Join(root, "target")
	mkdir(mainDir)
	writeAt(filepath.Join(root, "target-sibling-old"), "sibling sentinel", old)
	oldMatches, err := filepath.Glob(filepath.Join(root, "target*"))
	if err != nil || !slices.Contains(oldMatches, filepath.Join(root, "target-sibling-old")) {
		t.Fatalf("historical glob did not demonstrate sibling selection: %q, %v", oldMatches, err)
	}
	mainDeleted := []string{}
	for _, name := range []string{"old", "equal-time-a", "equal-time-b", "exact-48h", "after-48h", "space name", "tab\tname", "line\nname", "-leading", "-", "quote'\"back\\slash", "dollar$star*question?bracket[", "$(touch injected)"} {
		at := old
		if name == "exact-48h" {
			at = now.Add(-48 * time.Hour)
		}
		if name == "after-48h" {
			at = now.Add(-48*time.Hour - time.Minute)
		}
		writeAt(filepath.Join(mainDir, name), name, at)
		mainDeleted = append(mainDeleted, filepath.Join("target", name))
	}
	writeAt(filepath.Join(mainDir, "before-48h"), "survive", now.Add(-48*time.Hour+time.Minute))
	writeAt(filepath.Join(mainDir, "new"), "survive", newer)
	writeAt(filepath.Join(mainDir, ".hidden-old"), "survive", old)
	mkdir(filepath.Join(mainDir, ".hidden-dir"))
	writeAt(filepath.Join(mainDir, ".hidden-dir", "old"), "survive", old)
	mkdir(filepath.Join(mainDir, "nested"))
	for _, name := range []string{"old", ".nested-hidden-old"} {
		writeAt(filepath.Join(mainDir, "nested", name), "delete", old)
		mainDeleted = append(mainDeleted, filepath.Join("target", "nested", name))
	}
	writeAt(filepath.Join(mainDir, "nested", "new"), "survive", newer)
	mkdir(filepath.Join(mainDir, "old-empty-dir"))
	setDirTime(filepath.Join(mainDir, "old-empty-dir"), old)
	writeAt(filepath.Join(root, "hardlink-sentinel"), "survive", old)
	if err := os.Link(filepath.Join(root, "hardlink-sentinel"), filepath.Join(mainDir, "hardlink-old")); err != nil {
		t.Fatal(err)
	}
	mainDeleted = append(mainDeleted, filepath.Join("target", "hardlink-old"))
	if err := os.Symlink(filepath.Join(root, "hardlink-sentinel"), filepath.Join(mainDir, "symlink-file")); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(filepath.Join(root, "missing-symlink-target"), filepath.Join(mainDir, "broken-symlink")); err != nil {
		t.Fatal(err)
	}
	run("older files relative no slash", mainDir, "target", "f", "+", "1", false, 0, mainDeleted, false)
	if _, err := os.Lstat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("command substitution pathname executed: %v", err)
	}

	newDir := filepath.Join(root, "newer")
	mkdir(newDir)
	newDeleted := []string{filepath.Join("newer", "one-hour"), filepath.Join("newer", "before-24h"), filepath.Join("newer", "exact-24h")}
	writeAt(filepath.Join(newDir, "one-hour"), "delete", newer)
	writeAt(filepath.Join(newDir, "before-24h"), "delete", now.Add(-24*time.Hour+time.Minute))
	writeAt(filepath.Join(newDir, "exact-24h"), "survive", now.Add(-24*time.Hour))
	writeAt(filepath.Join(newDir, "after-24h"), "survive", now.Add(-24*time.Hour-time.Minute))
	writeAt(filepath.Join(newDir, "old"), "survive", old)
	writeAt(filepath.Join(newDir, ".hidden-new"), "survive", newer)
	run("newer files absolute trailing slash under errexit", newDir, newDir+string(filepath.Separator), "f", "-", "1", true, 0, newDeleted, false)

	emptyDir := filepath.Join(root, "only-hidden")
	mkdir(emptyDir)
	writeAt(filepath.Join(emptyDir, ".hidden"), "survive", old)
	run("unmatched nonhidden glob", emptyDir, "only-hidden", "f", "+", "1", false, 0, []string{}, false)

	dirDir := filepath.Join(root, "directories")
	mkdir(dirDir)
	for _, name := range []string{"old-empty", "new-empty", ".hidden-empty", "old-nonempty"} {
		mkdir(filepath.Join(dirDir, name))
	}
	writeAt(filepath.Join(dirDir, "old-nonempty", "content"), "survive", newer)
	for _, name := range []string{"old-empty", ".hidden-empty", "old-nonempty"} {
		setDirTime(filepath.Join(dirDir, name), old)
	}
	setDirTime(filepath.Join(dirDir, "new-empty"), newer)
	run("directory mode partial failure", dirDir, "directories", "d", "+", "1", false, 1, []string{filepath.Join("directories", "old-empty")}, true)

	statusScopes := []struct {
		name       string
		failures   []string
		successes  []string
		wantStatus int
		diagnostic bool
	}{
		{"success-then-failure", []string{"z-nonempty"}, []string{"a-empty"}, 1, true},
		{"failure-then-success", []string{"a-nonempty"}, []string{"z-empty"}, 1, true},
		{"multiple-failures", []string{"a-nonempty", "z-nonempty"}, nil, 1, true},
		{"all-successes", nil, []string{"a-empty", "z-empty"}, 0, false},
	}
	for _, tc := range statusScopes {
		scope := filepath.Join(root, tc.name)
		mkdir(scope)
		for _, name := range tc.failures {
			path := filepath.Join(scope, name)
			mkdir(path)
			writeAt(filepath.Join(path, "content"), "survive", newer)
			setDirTime(path, old)
		}
		wantDeleted := make([]string, 0, len(tc.successes))
		for _, name := range tc.successes {
			path := filepath.Join(scope, name)
			mkdir(path)
			setDirTime(path, old)
			wantDeleted = append(wantDeleted, filepath.Join(tc.name, name))
		}
		run(tc.name, scope, tc.name, "d", "+", "1", false, tc.wantStatus, wantDeleted, tc.diagnostic)
	}

	depthDir := filepath.Join(root, "directory-depth")
	mkdir(depthDir)
	visibleParent := filepath.Join(depthDir, "visible-parent")
	mkdir(visibleParent)
	hiddenDescendant := filepath.Join(visibleParent, ".hidden-child")
	mkdir(hiddenDescendant)
	setDirTime(hiddenDescendant, newer)
	setDirTime(visibleParent, newer)
	hiddenImmediate := filepath.Join(depthDir, ".hidden-immediate")
	mkdir(hiddenImmediate)
	setDirTime(hiddenImmediate, newer)
	linkedOutside := filepath.Join(root, "linked-directory-sentinel")
	mkdir(linkedOutside)
	writeAt(filepath.Join(linkedOutside, "old"), "outside subtree", old)
	if err := os.Symlink(linkedOutside, filepath.Join(depthDir, "linked-directory")); err != nil {
		t.Fatal(err)
	}
	run("directory mode depth-first and symlink", depthDir, "directory-depth", "d", "-", "1", false, 0,
		[]string{filepath.Join("directory-depth", "visible-parent", ".hidden-child"), filepath.Join("directory-depth", "visible-parent")}, false)

	if os.Geteuid() != 0 {
		for _, errexit := range []bool{false, true} {
			name := "partial-normal"
			if errexit {
				name = "partial-errexit"
			}
			scope := filepath.Join(root, name)
			mkdir(scope)
			writeAt(filepath.Join(scope, "a-old"), "delete", old)
			locked := filepath.Join(scope, "z-locked")
			mkdir(locked)
			writeAt(filepath.Join(locked, "old"), "survive", old)
			if err := os.Chmod(locked, 0o500); err != nil {
				t.Fatal(err)
			}
			if err := os.Remove(filepath.Join(locked, "old")); !os.IsPermission(err) {
				if err == nil {
					t.Fatal("permission probe unexpectedly removed protected fixture")
				}
				t.Fatalf("permission probe: %v", err)
			}
			run("read-only directory partial failure "+name, scope, name, "f", "+", "1", errexit, 1, []string{filepath.Join(name, "a-old")}, true)
			if err := os.Chmod(locked, 0o700); err != nil {
				t.Fatal(err)
			}
		}
	}

	for _, name := range []string{"-directory", "-", "special space\ttab\nline'$\\*?[$(touch injected)"} {
		scope := filepath.Join(root, name)
		mkdir(scope)
		writeAt(filepath.Join(scope, "old"), "delete", old)
		run("odd search directory "+name, scope, name, "f", "+", "1", true, 0, []string{filepath.Join(name, "old")}, false)
	}
	linkedSearch := filepath.Join(root, "linked-search-real")
	mkdir(linkedSearch)
	writeAt(filepath.Join(linkedSearch, "old"), "delete", old)
	if err := os.Symlink(linkedSearch, filepath.Join(root, "linked-search-alias")); err != nil {
		t.Fatal(err)
	}
	run("trailing slash on search-directory symlink", linkedSearch, "linked-search-alias/", "f", "+", "1", false, 0,
		[]string{filepath.Join("linked-search-real", "old")}, false)

	reject := func(label, input string, errexit bool) {
		t.Helper()
		assertRoot()
		before := inventory()
		script := expand(`${pathToDirectory}`, "f", "+", "1")
		if errexit {
			script = "set -e\n" + script
		}
		script += `printf 'STATUS:%s\n' "$?"` + "\n"
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		cmd := exec.CommandContext(ctx, "sh", "-c", script)
		cmd.Dir = root
		cmd.Env = append(os.Environ(), "pathToDirectory="+input, "LC_ALL=C")
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if ctx.Err() != nil {
			t.Fatalf("%s timed out", label)
		}
		wantOut, wantExit := "STATUS:2\n", 0
		if errexit {
			wantOut, wantExit = "", 2
		}
		if exitCode(err) != wantExit || string(out) != wantOut || stderr.Len() == 0 {
			t.Fatalf("%s: exit=%d out=%q stderr=%q", label, exitCode(err), out, stderr.String())
		}
		if !reflect.DeepEqual(before, inventory()) {
			t.Fatalf("%s changed fixture", label)
		}
		data, err := os.ReadFile(outsidePath)
		if err != nil || string(data) != "outside sentinel" {
			t.Fatalf("%s changed outside sentinel", label)
		}
	}
	reject("empty search input", "", false)
	reject("missing search input under errexit", "missing-search", true)
	if _, err := os.Lstat(filepath.Join(root, "injected")); !os.IsNotExist(err) {
		t.Fatalf("command substitution in directory name executed: %v", err)
	}
}
