DESCRIPTION = "This is a simple http server."
LICENSE = "MIT"
BB_STRICT_CHECKSUM = "0"
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

inherit go

SRC_URI = "file://server.go"

do_compile() {
        ${GO} build -o ${S}/server ${S}/server.go 
}

do_install() {
	install -d "${D}/${bindir}"
	install  -m 0755 "${WORKDIR}/build/server" "${D}/${bindir}"

}

