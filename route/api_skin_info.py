import urllib.request

from .tool.func import *

async def api_skin_info(name = ''):
    with get_db_connect() as conn:
        name = skin_check(conn) if name == '' else './views/' + name + '/index.html'

        if not flask.request.args.get('all', None):
            json_address = re.sub(r"(((?!\.|\/).)+)\.html$", "info.json", name)
            try:
                json_data = json_loads(open(json_address, encoding='utf8').read())
            except:
                json_data = None

            if json_data:
                return flask.jsonify(json_data)
            else:
                return flask.jsonify({}), 404
        else:
            a_data = {}
            d_link_data = {
                "ACME" : "https://raw.githubusercontent.com/openNAMU/openNAMU-Skin-ACME/master/info.json",
                "Liberty" : "https://raw.githubusercontent.com/openNAMU/openNAMU-Skin-Liberty/master/info.json",
                "Before Namu" : "https://raw.githubusercontent.com/openNAMU/openNAMU-Skin-Before_Namu/master/info.json"
            }

            for i in load_skin(conn, skin_check(conn, 1), 1):
                json_address = re.sub(r"(((?!\.|\/).)+)\.html$", "info.json", './views/' + i + '/index.html')
                try:
                    json_data = json_loads(open(json_address, encoding='utf8').read())
                except:
                    json_data = None

                if json_data:
                    if i == skin_check(conn, 1):
                        json_data = {**json_data, **{ "main" : "true" }}

                    if "info_link" in json_data:
                        info_link = json_data["info_link"]
                    elif json_data["name"] in d_link_data:
                        info_link = d_link_data[json_data["name"]]
                    else:
                        info_link = 0

                    if info_link != 0:
                        try:
                            get_data = urllib.request.urlopen(info_link)
                        except:
                            get_data = None

                        if get_data and get_data.getcode() == 200:
                            try:
                                get_data = json_loads(get_data.read().decode())
                            except:
                                get_data = {}

                            if "skin_ver" in get_data:
                                json_data = {**json_data, **{ "lastest_version" : {
                                    "skin_ver" : get_data["skin_ver"]
                                }}}

                    a_data = {**a_data, **{ i : json_data }}

            return flask.jsonify(a_data)