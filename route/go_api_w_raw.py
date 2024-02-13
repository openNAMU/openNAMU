from .tool.func import *

def api_w_raw(db_set, name = 'Test', rev = '', exist_check = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()
        
        other_set = {}
        other_set["name"] = name
        other_set["rev"] = rev
        other_set["exist_check"] = exist_check
        other_set = json.dumps(other_set)

        if acl_check(name, 'render') != 1:
            if platform.system() == 'Linux':
                if platform.machine() in ["AMD64", "x86_64"]:
                    data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.bin"), db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
                else:
                    data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.bin"), db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
            else:
                if platform.machine() in ["AMD64", "x86_64"]:
                    data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".amd64.exe"), db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
                else:
                    data = subprocess.Popen([os.path.join(".", "route_go", "bin", sys._getframe().f_code.co_name + ".arm64.exe"), db_set, other_set], stdout = subprocess.PIPE).communicate()[0]

            data = data.decode('utf8')

            return flask.Response(response = data, status = 200, mimetype = 'application/json')
        else:
            return flask.jsonify({})