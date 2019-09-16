from .tool.func import *

def list_acl_2():
    
    
    div = '''
        <table id="main_table_set">
            <tbody>
                <tr>
                    <td id="main_table_width_quarter">''' + load_lang('document_name') + '''</td>
                    <td id="main_table_width_quarter">''' + load_lang('document_acl') + '''</td>
                    <td id="main_table_width_quarter">''' + load_lang('discussion_acl') + '''</td>
                    <td id="main_table_width_quarter">''' + load_lang('view_acl') + '''</td>
    '''
    
    sqlQuery("select title, decu, dis, view, why from acl where decu != '' or dis != '' or view != '' order by title desc")
    list_data = sqlQuery("fetchall")
    for data in list_data:
        if not re.search('^user:', data[0]) and not re.search('^file:', data[0]):
            acl = []
            for i in range(1, 4):
                if data[i] == 'admin':
                    acl += [load_lang('admin')]
                elif data[i] == 'user':
                    acl += [load_lang('member')]
                elif data[i] == '':
                    acl += [load_lang('normal')]
                else:
                    acl += [data[i]]

            div +=  '''
                <tr>
                    <td>
                        <a href="/w/''' + url_pas(data[0]) + '">' + data[0] + '''</a>
                    </td>
                    <td>''' + acl[0] + '''</td>
                    <td>''' + acl[1] + '''</td>
                    <td>''' + acl[2] + '''</td>
                </tr>
            '''
        
    div +=  '''
            </tbody>
        </table>
    '''
    
    return easy_minify(flask.render_template(skin_check(), 
        imp = [load_lang('acl_document_list'), wiki_set(), custom(), other2([0, 0])],
        data = div,
        menu = [['other', load_lang('return')]]
    ))