import re

def footnote(data, fol_num):
    a = 1
    tou = "<hr style='margin-top: 30px;' id='footnote'><div><br>"
    namu = []
    pop_re = re.compile('(?:\[\*([^\s]*)(?:\s((?:(?!\[|\]).)*))?\]|(\[각주\]))')
    while(1):
        b = pop_re.search(data)
        if(b):
            results = b.groups()
            try:
                if(not results[1] and results[0]):
                    i = 0
                    
                    while(1):
                        try:
                            if(namu[i] == results[0]):
                                none_this = 0
                                break
                            else:
                                i += 2
                        except:
                            none_this = 1
                            break
                            
                    if(none_this == 0):
                        data = pop_re.sub("<sup><a href='javascript:void(0);' onclick='folding(" + str(fol_num) + ");' id='rfn-" + str(a) + "'>[" + results[0] + "]</a></sup>" + \
                                        "<div class='popup' style='display: none;' id='folding_" + str(fol_num) + "'><a onclick='folding(" + str(fol_num) + ");'" + \
                                        " href='#fn-" + str(a) + "'>#d#" + results[0] + "#/d#</a> <a href='javascript:void(0);' onclick='folding(" + str(fol_num) + ");'>[X]</a> " + \
                                         namu[i + 1] + "</div>", data, 1)
                    else:
                        data = pop_re.sub("<sup><a href='javascript:void(0);' id='rfn-" + str(a) + "'>#d#" + results[0] + "#/d#</a></sup>", data, 1)
                else:
                    if(results[0]):                
                        namu += [results[0]]
                        namu += [results[1]]

                        tou += "<span id='footnote-list'><a href='#rfn-" + str(a) + "' id='fn-" + str(a) + "'>[" + results[0] + "]</a> " + results[1] + "</span><br>"
                        data = pop_re.sub("<sup><a href='javascript:void(0);' onclick='folding(" + str(fol_num) + ");' id='rfn-" + str(a) + "'>#d#" + results[0] + "#/d#</a>" + \
                                        "</sup><div class='popup' style='display: none;' id='folding_" + str(fol_num) + "'><a onclick='folding(" + str(fol_num) + ");'" + \
                                        " href='#fn-" + str(a) + "'>#d#" + results[0] + "#/d#</a> <a href='javascript:void(0);' onclick='folding(" + str(fol_num) + ");'>" + \
                                        "#d#X#/d#</a> " + results[1] + "</div>", data, 1)     
                    else:                    
                        tou += "<span id='footnote-list'><a href='#rfn-" + str(a) + "' id='fn-" + str(a) + "'>[" + str(a) + "]</a> " + results[1] + "</span><br>"
                        data = pop_re.sub('<sup><a href="javascript:void(0);" onclick="folding(' + str(fol_num) + ');" id="rfn-' + str(a) + '">#d#' + str(a) + '#/d#</a></sup>' + \
                                        '<div class="popup" style="display: none;" id="folding_' + str(fol_num) + '"><a onclick="folding(' + str(fol_num) + ');"' + \
                                        ' href="#fn-' + str(a) + '">#d#' + str(a) + '#/d#</a> <a href="javascript:void(0);" onclick="folding(' + str(fol_num) + ');">#d#X#/d#</a> ' + \
                                        results[1] + '</div>', data, 1)
                    a += 1

                fol_num += 2
            except:
                tou += '</div>'

                if(tou == "<hr style='margin-top: 30px;' id='footnote'><div><br></div>"):
                    tou = ""
                else:
                    tou = re.sub('#d#(?P<in>(?:(?!#\/d#).)*)#\/d#', '[\g<in>]', tou)

                data = pop_re.sub("<br>" + tou, data, 1)
                tou = "<hr style='margin-top: 30px;' id='footnote'><div><br>"
        else:
            tou += '</div>'

            if(tou == "<hr style='margin-top: 30px;' id='footnote'><div><br></div>"):
                tou = ""
            else:
                tou = re.sub('#d#(?P<in>(?:(?!#\/d#).)*)#\/d#', '[\g<in>]', tou)

            break
            
    data = re.sub('#d#(?P<in>(?:(?!#\/d#).)*)#\/d#', '[\g<in>]', data)
    
    data = re.sub("\[각주\](?:(?:<br>| |\r|\n)+)?$", "", data)
    data = re.sub("(?:(?:<br>| |\r|\n)+)$", "", data)
    data += tou
    
    return(data)