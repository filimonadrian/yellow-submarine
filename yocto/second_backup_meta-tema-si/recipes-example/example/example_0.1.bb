SUMMARY = "bitbake-layers recipe"
DESCRIPTION = "Recipe created by bitbake-layers"
LICENSE = "MIT"
BB_STRICT_CHECKSUM = "0"
# 2. fișierul licenței alese (MIT) e partajat și are acest checksum:
LIC_FILES_CHKSUM = "file://${COMMON_LICENSE_DIR}/MIT;md5=0835ade698e0bcf8506ecda2f7b4f302"

# adăugăm fișiere sursă
SRC_URI = "file://hello.py"

RDEPENDS_${PN} = "python3-core"

# fișierele sursă vor fi copiate / descărcate în ${WORKDIR}! 
# va trebui să intervenim cu un pas de instalare pentru a le copia în imagine:
do_install() {
    # atenție: ${D} represintă directorul destinație în procesul de împachetare!
    #          orice generați în afara acestuia nu va fi inclus în pachet!
    # copiem scriptul în bin:
    install -D -m 0755 -t ${D}${bindir} ${WORKDIR}/hello.py
}
