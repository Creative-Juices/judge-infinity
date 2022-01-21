#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#define LIMIT 1000000
#define end_ch(chvar) (chvar == EOF || chvar == '\n' || chvar == '\r')
#define sp_ch(chvar) (chvar == ' ' || chvar == '\t')

FILE *f1, *f2;
char line1[LIMIT], line2[LIMIT];

void clearLineStrings() {
    memset( line1, '\0', sizeof(char)*LIMIT );
    memset( line2, '\0', sizeof(char)*LIMIT );
}

int compare(char line1[], char line2[]) {
    int p1 = 0, p2 = 0, n1 = strlen(line1) - 1, n2 = strlen(line2) - 1;
    while (sp_ch(line1[p1])) p1++;
    while (sp_ch(line2[p2])) p2++;
    while (sp_ch(line1[n1]) || end_ch(line1[n1])) n1--;
    while (sp_ch(line2[n2]) || end_ch(line2[n2])) n2--;
    for (; p1 <= n1 && p2 <= n2; p1++, p2++) if (line1[p1] != line2[p2]) return 0;
    return (p1 <= n1 || p2 <= n2) ? 0 : 1;
}

int main(int argc, const char *argv[]) {
    if( argc < 3 ) exit(3); // Missing Arguments

    f1 = fopen(argv[1], "r");
    f2 = fopen(argv[2], "r");

    while (feof(f1) == 0 && feof(f2) == 0) { // Do not put fgets directly here, one of them might get read and reach eof and might indicate AC
        clearLineStrings();
        fgets(line1, LIMIT, f1);
        fgets(line2, LIMIT, f2);
        if (compare(line1, line2) == 0) exit(4); // Wrong Answer
    }

    int f1ended = feof(f1), f2ended = feof(f2);
    if( f1ended && f2ended ) exit(0); // Success
    else if ( f1ended ) {
        while( feof(f2) ){
            clearLineStrings();
            fgets(line2, LIMIT, f2);
            for( int i = 0 ; line2[i] != EOF && line2[i] != '\0' ; i++ ) {
                if( !end_ch(line2[i]) && !sp_ch(line2[i]) ) exit(4); // Wrong Answer
            }
        }
    }else{
        while( feof(f1) ){
            clearLineStrings();
            fgets(line1, LIMIT, f1);
            for( int i = 0 ; line1[i] != EOF && line1[i] != '\0' ; i++ ) {
                if( !end_ch(line1[i]) && !sp_ch(line1[i]) ) exit(4); // Wrong Answer
            }
        }
    }

    exit(0); // Success
}