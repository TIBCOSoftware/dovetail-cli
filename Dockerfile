###################################################################################
#                          buildtype/dovetail                                        #
###################################################################################
FROM scratch

VOLUME [ /var/lib/build_server/buildtypes/dovetail ]

COPY . /var/lib/build_server/buildtypes/dovetail/