from .tool.func import *

def setting_adsense_2(conn):
    curs = conn.cursor()

    if admin_check(None, 'adsense setting') != 1:
        return re_error('/error/3')
    
    if flask.request.method == 'POST':
        try:
            adsense_enabled = flask.request.form.get('adsense_enabled')
            adsense_code = flask.request.form['adsense_code']
        except:
            return easy_minify(flask.render_template(skin_check(),
                imp = [load_lang('inter_error'), wiki_set(), custom(), other2([0, 0])],
                data = '<h2>ie_no_data_required</h2>' + load_lang('ie_no_data_required'),
                menu = [['other', load_lang('return')]]
            ))
        
        if adsense_enabled == 'on':
            curs.execute('update other set data = "True" where name = "adsense"')
        else:
            curs.execute('update other set data = "False" where name = "adsense"')
        
        curs.execute('update other set data = ? where name = "adsense_code"', [adsense_code])
        conn.commit()
        
        return redirect('/adsense_setting')

    body_content = ''

    curs.execute('select data from other where name = "adsense"')
    adsense_enabled = curs.fetchall()[0][0]

    curs.execute('select data from other where name = "adsense_code"')
    adsense_code = curs.fetchall()[0][0]

    template = '''
        <form action="" accept-charset="utf-8" method="post">
            <div class="form-check">
                <label class="form-check-label">
                    <input class="form-check-input" name="adsense_enabled" type="checkbox" {}>
                    {}
                </label>
            </div>
            <hr>
            <div class="form-group">
                <textarea class="form-control" id="adsense_code" name="adsense_code" rows="12">{}</textarea>
            </div>
            <button type="submit" value="publish">{}</button>
        </form>
    '''
    
    body_content += template.format(
        'checked' if adsense_enabled == 'True' else template.format(''),
        load_lang('adsense_enable'),
        load_lang('save'),
        adsense_code
    )

    return easy_minify(flask.render_template(skin_check(),
        imp = [load_lang('adsense_setting'), wiki_set(), custom(), other2([0, 0])],
        data = body_content,
        menu = [['other', load_lang('return')]]
    ))