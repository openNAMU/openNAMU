from .tool.func import *

def api_w_xref(db_set, name = 'Test', page = 1, xref_type = '1'):
    with get_db_connect() as conn:
        other_set = {}
        other_set["name"] = name
        other_set["page"] = str(page)
        other_set["do_type"] = xref_type
        other_set = json.dumps(other_set)

        if platform.system() == 'Linux':
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.bin"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.bin"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            if platform.machine() in ["AMD64", "x86_64"]:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.exe"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]
            else:
                data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.exe"), sys._getframe().f_code.co_name, db_set, other_set], stdout = subprocess.PIPE).communicate()[0]

        data = data.decode('utf8')

        return flask.Response(response = data, status = 200, mimetype = 'application/json')