from .tool.func import *

def api_list_recent_change(num = 1, set_type = 'normal', limit = 10, legacy = 'on'):
    other_set = {}
    other_set["num"] = str(num)
    other_set["limit"] = str(limit)
    other_set["set_type"] = set_type
    other_set["legacy"] = legacy
    other_set["ip"] = ip_check()
    other_set = json.dumps(other_set)

    if platform.system() == 'Linux':
        if platform.machine() in ["AMD64", "x86_64"]:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.bin"), sys._getframe().f_code.co_name, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.bin"), sys._getframe().f_code.co_name, other_set], stdout = subprocess.PIPE).communicate()[0]
    else:
        if platform.machine() in ["AMD64", "x86_64"]:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.amd64.exe"), sys._getframe().f_code.co_name, other_set], stdout = subprocess.PIPE).communicate()[0]
        else:
            data = subprocess.Popen([os.path.join(".", "route_go", "bin", "main.arm64.exe"), sys._getframe().f_code.co_name, other_set], stdout = subprocess.PIPE).communicate()[0]

    data = data.decode('utf8')

    response = flask.Response(response = data, status = 200, mimetype = 'application/json')
    response.headers.add("Access-Control-Allow-Origin", "*")
    response.headers.add('Access-Control-Allow-Headers', "Content-Type")
    response.headers.add('Access-Control-Allow-Methods', "GET")

    return response