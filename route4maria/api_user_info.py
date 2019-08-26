from .tool.func import *
import pymysql

def api_user_info_2(conn, name):
    curs = conn.cursor()

    if flask.request.args.get('render', None):
        plus_d = ''
        plus_t = []

        plus_d = '''
            <table>
                <tbody>
                    <tr>
                        <td>''' + load_lang('user_name') + '''</td>
                        <td>{}</td>
                    </tr>
                    <tr>
                        <td>''' + load_lang('authority') + '''</td>
                        <td>{}</td>
                    </tr>
                    <tr>
                        <td>''' + load_lang('state') + '''</td>
                        <td>{}</td>
                    </tr>
                </tbody>
            </table>
        '''

        curs.execute("select acl from user where id = %s", [name])
        data = curs.fetchall()
        if data:
            if data[0][0] != 'user':
                plus_t += [data[0][0]]
            else:
                plus_t += [load_lang('member')]
        else:
            plus_t += [load_lang('normal')]

        if ban_check(name) == 0:
            plus_t += [load_lang('normal')]
        else:
            match = re.search("^([0-9]{1,3}\.[0-9]{1,3})", name)
            if match:
                match = match.groups()[0]
            else:
                match = '-'

            curs.execute("select end, login, band from ban where block = %s or block = %s", [name, match])
            block_data = curs.fetchall()
            if block_data:
                if block_data[0][0] != '':
                    plus_t += [load_lang('period') + ' : ' + block_data[0][0]]
                else:
                    plus_t += [load_lang('limitless')]

                if block_data[0][1] != '':
                    plus_t += [load_lang('login_able')]

                if block_data[0][2] == 'O':
                    plus_t += [load_lang('band_blocked')]

        plus_d = plus_d.format(ip_pas(name), plus_t[0], plus_t[1])

        return flask.jsonify({ "data" : plus_d })
    else:
        pass 