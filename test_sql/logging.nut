#ifndef _LOGGING_NUT
#define _LOGGING_NUT

/**
 Logging Library
 */
 
#mode string QQQ "$$" "$$"
#define __LINE__ #line
#define __FILE__ #file
#define log(txt) dbg($$"__FILE__:__LINE__"$$, txt)

function dbg(fileline, statement) {
}

#endif 

