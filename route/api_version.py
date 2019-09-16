from .tool.func import *

def api_version_2(r_ver, c_ver):
    

    n_ver = ''
    data = None

    data = urllib.request.urlopen('https://raw.githubusercontent.com/2du/openNAMU/master/version.json')
    if data and data.getcode() == 200:
        try:
            json_data = json.loads(data.read().decode())
            if 'master' in json_data:
                n_ver = json_data['master']['r_ver']
        except:
            pass
        
    json_data = { "version" : r_ver, "db_version" : c_ver, "lastest_version" : n_ver  }

    return flask.jsonify(json_data)