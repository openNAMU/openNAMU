from .tool.func import *

def api_user_info(name = ''):
    with get_db_connect() as conn:
        curs = conn.cursor()

        if flask.request.method == 'POST':
            try:
                data_list = json.loads(flask.request.form.get('title_list', ''))
                data_list = list(set(data_list))
            except:
                data_list = [name]

            data_result = {}
            for user_name in data_list:
                data_result[user_name] = {}
                
                # name part
                data_result[user_name]['render'] = ip_pas(user_name)
                
                # auth part
                curs.execute(db_change("select data from user_set where id = ? and name = 'acl'"), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    if db_data[0][0] != 'user':
                        curs.execute(db_change("select name from alist where name = ?"), [db_data[0][0]])
                        if curs.fetchall() or db_data[0][0] in get_default_admin_group():
                            data_result[user_name]['auth'] = db_data[0][0]
                        else:
                            data_result[user_name]['auth'] = '1'
                    else:
                        data_result[user_name]['auth'] = '1'
                else:
                    data_result[user_name]['auth'] = '0'

                curs.execute(db_change("select data from user_set where id = ? and name = 'auth_date'"), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    data_result[user_name]['auth_date'] = db_data[0][0]
                else:
                    data_result[user_name]['auth_date'] = '0'

                curs.execute(db_change("select data from user_set where id = ? and name = 'level'"), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    data_result[user_name]['level'] = db_data[0][0]
                else:
                    data_result[user_name]['level'] = '0'

                curs.execute(db_change("select data from user_set where id = ? and name = 'experience'"), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    data_result[user_name]['exp'] = db_data[0][0]
                else:
                    data_result[user_name]['exp'] = '0'
                    
                # ban part
                if ban_check(name) == 0:
                    data_result[user_name]['ban'] = '0'
                else:
                    data_result[user_name]['ban'] = {}
                    regex_ban = 0
                    
                    curs.execute(db_change("select login, block, end, why from rb where band = 'regex' and ongoing = '1'"))
                    for db_data in curs.fetchall():
                        if re.compile(db_data[1]).search(user_name):
                            regex_ban = 1
                            
                            data_result[user_name]['ban']['type'] = 'regex'
                            if db_data[0] == 'O':
                                data_result[user_name]['ban']['login_able'] = '1'
                            else:
                                data_result[user_name]['ban']['login_able'] = '0'
                                
                            if db_data[2] == '':
                                data_result[user_name]['ban']['period'] = '0'
                            else:
                                data_result[user_name]['ban']['period'] = db_data[2]
                                
                            data_result[user_name]['ban']['reason'] = db_data[3]
                            
                            break
                            
                    if regex_ban == 0:
                        curs.execute(db_change("select login, block, end, why from rb where block = ? and ongoing = '1'"), [user_name])
                        db_data = curs.fetchall()
                        if db_data:
                            data_result[user_name]['ban']['type'] = 'normal'
                            if db_data[0][0] == 'O':
                                data_result[user_name]['ban']['login_able'] = '1'
                            else:
                                data_result[user_name]['ban']['login_able'] = '0'
                                
                            if db_data[0][2] == '':
                                data_result[user_name]['ban']['period'] = '0'
                            else:
                                data_result[user_name]['ban']['period'] = db_data[0][2]
                                
                            data_result[user_name]['ban']['reason'] = db_data[0][3]
                
                # user document part
                curs.execute(db_change("select title from data where title = ?"), ['user:' + user_name])
                if curs.fetchall():
                    data_result[user_name]['document'] = '1'
                else:
                    data_result[user_name]['document'] = '0'

                # user title part
                curs.execute(db_change('select data from user_set where name = "user_title" and id = ?'), [user_name])
                db_data = curs.fetchall()
                if db_data:
                    data_result[user_name]['user_title'] = db_data[0][0]
                else:
                    data_result[user_name]['user_title'] = ''
                    
            return flask.jsonify(data_result)
        else:
            return flask.jsonify({})