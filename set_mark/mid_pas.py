import re

def mid_pas(data, fol_num, include, in_c, toc_y):
    p = re.compile('{{{((?:(?:(?:\+|-)[0-5])|(?:#|@)(?:(?:[0-9a-f-A-F]{3}){1,2}|(?:\w+))|(?:#!(?:html|wiki|noin|folding|syntax)))(?:(?!{{{|}}}).)+)}}}', re.DOTALL)
    while 1:
        m = p.search(data)
        if m:
            data = p.sub('###' + m.groups()[0] + '/###', data, 1)
        else:
            break

    com = re.compile("{{{((?:(?!{{{|}}}).)*)}}}", re.DOTALL)
    while 1:
        m = com.search(data)
        if m:
            data = com.sub('<code>' + m.groups()[0] + '</code>', data, 1)
        else:
            break

    com3 = re.compile('###((?:(?!\/###).)+)\/###', re.DOTALL)
    m = com3.search(data)
    while 1:
        m = com3.search(data)
        if m:
            data = com3.sub('{{{' + m.groups()[0] + '}}}', data, 1)
        else:
            break

    com2 = re.compile("<code>((?:(?!(?:<code>|<\/code>)).)*)<\/code>", re.DOTALL)
    da_com = com2.findall(data)
    for com_da in da_com:
        mid_data = com_da.replace('<', '&lt;').replace('>', '&gt;')
        mid_data = re.sub("(?P<in>.)", "#no#\g<in>#/no#", mid_data)
        data = com2.sub(mid_data, data, 1)

    while 1:
        is_it = com.search(data)
        if is_it:
            it_d = is_it.groups()[0]

            big_a = re.compile("^\+([1-5]) (.*)$", re.DOTALL)
            big = big_a.search(it_d)

            small_a = re.compile("^\-([1-5]) (.*)$", re.DOTALL)
            small = small_a.search(it_d)

            color_b = re.compile("^(#(?:[0-9a-f-A-F]{3}){1,2}) (.*)$", re.DOTALL)
            color_2 = color_b.search(it_d)

            color_c = re.compile("^#(\w+) (.*)$", re.DOTALL)
            color_3 = color_c.search(it_d)

            back_a = re.compile("^@((?:[0-9a-f-A-F]{3}){1,2}) (.*)$", re.DOTALL)
            back = back_a.search(it_d)

            back_c = re.compile("^@(\w+) (.*)$", re.DOTALL)
            back_3 = back_c.search(it_d)

            include_out_a = re.compile("^#!noin (.*)$", re.DOTALL)
            include_out = include_out_a.search(it_d)

            div_a = re.compile("^#!wiki style=(?:&quot;|&#x27;)((?:(?!&quot;|&#x27;).)*)(?:&quot;|&#x27;)\r\n(.*)$", re.DOTALL)
            div = div_a.search(it_d)

            html_a = re.compile("^#!html (.*)$", re.DOTALL)
            html_d = html_a.search(it_d)

            fol_a = re.compile("^#!folding ((?:(?!\n).)*)\n(.*)$", re.DOTALL)
            fol = fol_a.search(it_d)

            syn_a = re.compile("^#!syntax ((?:(?!\n).)*)\n(.*)$", re.DOTALL)
            syn = syn_a.search(it_d)

            if big:
                big_d = big.groups()
                data = com.sub('<span style="font-size: ' + str(int(big_d[0]) * 20 + 100) + '%;">' + big_d[1] + '</span>', data, 1)
            elif small:
                sm_d = small.groups()
                data = com.sub('<span style="font-size: ' + str(100 - int(sm_d[0]) * 10) + '%;">' + sm_d[1] + '</span>', data, 1)
            elif color_2:
                c_d_2 = color_2.groups()
                data = com.sub('<span style="color: ' + c_d_2[0] + '">' + c_d_2[1] + '</span>', data, 1)
            elif color_3:
                c_d_3 = color_3.groups()
                data = com.sub('<span style="color: ' + c_d_3[0] + '">' + c_d_3[1] + '</span>', data, 1)
            elif back:
                back_d_1 = back.groups()
                data = com.sub('<span style="background: #' + back_d_1[0] + '">' + back_d_1[1] + '</span>', data, 1)
            elif back_3:
                back_d_3 = back_3.groups()
                data = com.sub('<span style="background: ' + back_d_3[0] + '">' + back_d_3[1] + '</span>', data, 1)
            elif div:
                div_d = div.groups()
                data = com.sub('<div style="' + div_d[0] + '">' + div_d[1] + '</div>', data, 1)
            elif html_d:
                data = com.sub(html_d.groups()[0], data, 1)
            elif fol:
                fol_d = fol.groups()
                if toc_y != 0:
                    data = com.sub("<div>" + fol_d[0] + " <div id='folding_" + str(fol_num + 1) + "' style='display: inline-block;'>[<a href='javascript:void(0);' onclick='folding(" + str(fol_num + 1) + \
                                    "); folding(" + str(fol_num + 2) + "); folding(" + str(fol_num) + ");'>펼치기</a>]</div><div id='folding_" + str(fol_num + 2) + \
                                    "' style='display: none;'>[<a href='javascript:void(0);' onclick='folding(" + str(fol_num + 1) + "); folding(" + str(fol_num + 2) + \
                                    "); folding(" + str(fol_num) + ");'>접기</a>]</div><div id='folding_" + str(fol_num) + "' style='display: none;'><br>" + fol_d[1] + \
                                    "</div></div>", data, 1)
                    fol_num += 3
                else:
                    data = com.sub("<div>" + fol_d[0] + "<br><br>" + fol_d[1] + "</div>", data, 1)
                
            elif syn:
                syn_d = syn.groups()
                tax_d = syn_d[1].replace(' ', '<space>')
                tax_d = tax_d.replace('\r\n', '<isbr>')
                data = com.sub('<pre id="syntax"><code class="' + syn_d[0] + '"><code>' + tax_d + '</code></code></pre>', data, 1)
            elif include_out:
                if (include or in_c) == 1:
                    data = com.sub("", data, 1)
                else:
                    data = com.sub(include_out.groups()[0], data, 1)
            else:
                data = com.sub(it_d, data, 1)
        else:
            break

    com2 = re.compile("<code>((?:(?!(?:<code>|<\/code>)).)*)<\/code>", re.DOTALL)
    da_com = com2.findall(data)
    for com_da in da_com:
        mid_data = com_da.replace('<', '&lt;').replace('>', '&gt;')
        mid_data = re.sub("(?P<in>.)", "#no#\g<in>#/no#", mid_data)
        data = com2.sub(mid_data, data, 1)
            
    return [data, fol_num]