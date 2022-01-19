
IMAGE_INSTALL += " dhcp-client go avahi-daemon ssh server gui"

inherit extrausers

EXTRA_USERS_PARAMS += "usermod -P labsi root;"
EXTRA_USERS_PARAMS += "useradd -P student stud;"


