from .tool.func import *

def server_indexing_2():
    

    if admin_check() != 1:
        return re_error('/error/3')

    if flask.request.method == 'POST':
        admin_check(None, 'indexing')

        sqlQuery("select name from sqlite_master where type = 'index'")
        data = sqlQuery("fetchall")
        if data:
            for delete_index in data:
                print('Delete : ' + delete_index[0])

                sql = 'drop index if exists ' + delete_index[0]
                
                try:
                    sqlQuery(sql)
                except:
                    pass
        else:
            sqlQuery("select name from sqlite_master where type in ('table', 'view') and name not like 'sqlite_%' union all select name from sqlite_temp_master where type in ('table', 'view') order by 1;")
            for table in sqlQuery("fetchall"):            
                sqlQuery('select sql from sqlite_master where name = ?', [table[0]])
                cul = sqlQuery("fetchall")
                
                r_cul = re.findall('(?:([^ (]*) text)', str(cul[0]))
                
                for n_cul in r_cul:
                    print('Create : index_' + table[0] + '_' + n_cul)

                    sql = 'create index index_' + table[0] + '_' + n_cul + ' on ' + table[0] + '(' + n_cul + ')'
                    try:
                        sqlQuery(sql)
                    except:
                        pass

        sqlQuery("commit")
        
        return redirect()  
    else:
        sqlQuery("select name from sqlite_master where type = 'index'")
        data = sqlQuery("fetchall")
        if data:
            b_data = load_lang('delete')
        else:
            b_data = load_lang('create')

        return easy_minify(flask.render_template(skin_check(), 
            imp = [load_lang('indexing'), wiki_set(), custom(), other2([0, 0])],
            data =  '''
                    <form method="post">
                        <button type="submit">''' + b_data + '''</button>
                    </form>
                    ''',
            menu = [['manager', load_lang('return')]]
        ))   