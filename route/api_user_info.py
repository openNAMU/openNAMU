from .tool.func import *

async def api_user_info(user_name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        data_result = {}
        
        # name part
        data_result['render'] = await ip_pas(user_name)
        
        # auth part
        curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [user_name])
        db_data = curs.fetchall()
        if db_data:
            data_result['auth'] = db_data[0][0]
        elif ip_or_user(user_name) == 1:
            data_result['auth'] = 'ip'
        else:
            data_result['auth'] = 'user'

        curs.execute(db_change("select data from user_set where id = ? and name = 'auth_date'"), [user_name])
        db_data = curs.fetchall()
        if db_data:
            data_result['auth_date'] = db_data[0][0]
        else:
            data_result['auth_date'] = '0'

        level_data = await level_check(user_name)
        data_result['level'] = level_data[0]
        data_result['exp'] = level_data[1]
        data_result['max_exp'] = level_data[2]
            
        # ban part
        ban = await ban_check(user_name)
        if ban[0] == 0:
            data_result['ban'] = '0'
        else:
            data_result['ban'] = ban
        
        # user document part
        curs.execute(db_change("select title from data where title = ?"), ['user:' + user_name])
        if curs.fetchall():
            data_result['document'] = '1'
        else:
            data_result['document'] = '0'

        # user title part
        curs.execute(db_change('select data from user_set where name = "user_title" and id = ?'), [user_name])
        db_data = curs.fetchall()
        if db_data:
            data_result['user_title'] = db_data[0][0]
        else:
            data_result['user_title'] = ''

        lang_data_list = [
            'user_name',
            'authority',
            'state',
            'member',
            'normal',
            'blocked',
            'type',
            'regex',
            'period',
            'limitless',
            'login_able',
            'why',
            'band_blocked',
            'ip',
            'ban',
            'level',
            'option',
            'edit_request_able',
            'cidr'
        ]
        lang_data = { for_a : get_lang(conn, for_a) for for_a in lang_data_list }
                
        return flask.jsonify({ 'data' : data_result, 'language' : lang_data })