import os

env_dict = {
    'host' : os.getenv('NAMU_HOST'),
    'port' : os.getenv('NAMU_PORT'),
    'language' : os.getenv('NAMU_LANG'),
    'markup' : os.getenv('NAMU_MARKUP'),
    'encode' : os.getenv('NAMU_ENCRYPT')
}

server_set_var = {
    'host' : {
        'display' : 'Host',
        'require' : 'conv',
        'default' : '0.0.0.0'
    },
    'port' : {
        'display' : 'Port',
        'require' : 'conv',
        'default' : '3000'
    },
    'language' : {
        'display' : 'Language',
        'require' : 'select',
        'default' : 'en-US',
        'list' : ['ko-KR', 'en-US']
    },
    'markup' : {
        'display' : 'Markup',
        'require' : 'select',
        'default' : 'namumark',
        'list' : ['namumark', 'markdown', 'raw']
    },
    'encode' : {
        'display' : 'Encryption method',
        'require' : 'select',
        'default' : 'sha3',
        'list' : ['sha3', 'sha256']
    }
}

def init(key):
    if env_dict[key] != None:
        return env_dict[key]
    else:
        if key == 'markup':
            return server_set_var[key]['default']

        while 1:
            if server_set_var[key]['require'] == 'select':
                list_ = '[' + ', '.join(server_set_var[key]['list']) + ']'
            else:
                list_ = ''

            print('{} ({}) {} : '.format(
                server_set_var[key]['display'],
                server_set_var[key]['default'],
                list_
            ), end = '')

            server_set_val = input()
            if server_set_val:
                if server_set_var[key]['require'] == 'select':
                    if server_set_val not in server_set_var[key]['list']:
                        pass
                    else:
                        return server_set_val
                else:
                    return server_set_val
            else:
                return server_set_var[key]['default']