
IMAGE_INSTALL += " dhcp-client go avahi-daemon ssh go-server"

inherit extrausers

EXTRA_USERS_PARAMS += "usermod -P labsi root;"
EXTRA_USERS_PARAMS += "useradd -P student student;"


