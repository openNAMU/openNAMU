import re
import html

def html_pas(data):
    data = re.sub('%H%', '<', data)
    data = re.sub('%\/H%', '>', data)

    d_list = re.findall('<(\/)?([^> ]+)( (?:[^>]+)?)?>', data)
    for i_list in d_list:
        if i_list[0] == '':
            if i_list[1] in ['div', 'span', 'embed', 'iframe', 'ruby', 'rp', 'rt']:
                if re.search('<\/' + i_list[1] + '>', data):
                    src = re.search('src=([^ ]*)', i_list[2])
                    if src:
                        v_src = re.search('http(?:s)?:\/\/([^/\'" ]*)', src.groups()[0])
                        data_list = ["www.youtube.com", "serviceapi.nmv.naver.com", "tv.kakao.com", "www.google.com", "serviceapi.rmcnmv.naver.com"]
                        if v_src:
                            if not v_src.groups()[0] in data_list:
                                ot = re.sub('src=([^ ]*)', '', i_list[2])
                            else:
                                ot = i_list[2]
                        else:
                            ot = re.sub('src=([^ ]*)', '', i_list[2])
                    else:
                        ot = i_list[2]

                    po = re.compile('position', re.I)
                    data = data.replace('<' + i_list[1] + i_list[2] + '>', '%H%' + i_list[1] + po.sub('', ot) + '%/H%', 1)
                    data = re.sub('<\/' + i_list[1] + '>', '%H%/' + i_list[1] + '%/H%', data, 1)

    data = html.escape(data)
    data = data.replace('\\', '&#92;')

    print([data])
    
    end = re.findall('%H%((?:(?!%/H%).)*)%/H%', data)
    for d_end in end:
        data = re.sub('%H%((?:(?!%/H%).)*)%/H%', '<' + re.sub('&quot;', '"', re.sub('&#x27;', "'", d_end)) + '>', data, 1)

    return data