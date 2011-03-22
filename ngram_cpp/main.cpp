#include <QtCore/QCoreApplication>
#include <iostream>
#include <QFile>
#include <QByteArray>
#include <QHash>
#include <stdlib.h>

using namespace std;


// Ngram1
class ngram1 {
public:
        QByteArray data;
};
inline bool operator==(const ngram1 &e1, const ngram1 &e2)
{
    return e1.data == e2.data;
}
inline uint qHash(const ngram1& key)
{
    return qHash(key.data);
}

//Ngram2
class ngram2{
public:
        QByteArray data[2];
};
inline bool operator==(const ngram2 &e1, const ngram2 &e2)
{
    return e1.data[0] == e2.data[0]
            && e1.data[1] == e2.data[1];
}
inline uint qHash(const ngram2& key)
{
    return qHash(QString(key.data[0]) + " " + QString(key.data[1]));
}

//Ngram3
class ngram3 {
public:
        QByteArray data[3];
};
inline bool operator==(const ngram3 &e1, const ngram3 &e2)
{
    return e1.data[0] == e2.data[0]
            && e1.data[1] == e2.data[1]
            && e1.data[2] == e2.data[2];
}
inline uint qHash(const ngram3& key)
{
    return qHash(QString(key.data[0]) + " " +
                 QString(key.data[1]) + " " +
                 QString(key.data[2]));
}

//Ngram4
class ngram4 {
public:
    QByteArray data[4];
};
inline bool operator==(const ngram4 &e1, const ngram4 &e2)
{
    return e1.data[0] == e2.data[0]
            && e1.data[1] == e2.data[1]
            && e1.data[2] == e2.data[2]
            && e1.data[3] == e2.data[3];
}
inline uint qHash(const ngram4& key)
{
    return qHash(QString(key.data[0]) + " " +
                 QString(key.data[1]) + " " +
                 QString(key.data[2]) + " " +
                 QString(key.data[3]));
}

// maps to hold the Ngrams
QHash<ngram1, quint64> m1;
QHash<ngram2, quint64> m2;
QHash<ngram3, quint64> m3;
QHash<ngram4, quint64> m4;

//storing ngrams in maps

ngram1 ng1;
void storeNgram1 (QByteArray word) {
    ng1.data = word;
    if (m1.contains(ng1) == false)
    {
        m1[ng1] = 1;
    }
    else {
        m1[ng1]++;
    }
}

ngram2 ng2;
void storeNgram2 (QByteArray word) {
    ng2.data[0] = ng2.data[1];
    ng2.data[1] = word;
    QByteArray temp = ng2.data[0];
    if (ng2.data[1].length()>0)
    {
        if (m2.contains(ng2) == false)
        {
            m2[ng2] = 1;
        }
        else {
            m2[ng2]++;
        }
    }
}

ngram3 ng3;
void storeNgram3 (QByteArray word) {
    ng3.data[0] =  ng3.data[1];
    ng3.data[1] =  ng3.data[2];
    ng3.data[2] = word; 
    if (ng3.data[1].length()>0 && ng3.data[0].length()>0)
    {
        if (m3.contains(ng3) == false)
        {
            m3[ng3] = 1;
        }
        else {
            m3[ng3]++;
        }
    }
    //cout << ng3.data[2].data() << " " << ng3.data[1].data() << " " <<ng3.data[0].data() << endl;
    //cout << "val = " << m3[ng3] << endl;
}

ngram4 ng4;
void storeNgram4 (QByteArray word) {
    ng4.data[0] =  ng4.data[1];
    ng4.data[1] =  ng4.data[2];
    ng4.data[2] =  ng4.data[3];
    ng4.data[3] = word;
    if (ng4.data[2].length()>0 && ng4.data[1].length()>0 && ng4.data[0].length()>0)
    {
        if (m4.contains(ng4) == false)
        {
            m4[ng4] = 1;
        }
        else {
            m4[ng4]++;
        }
    }
}



int main(int argc, char *argv[])
{

    cout << "starting" << endl;
    uint linenumber = 0;
    uint wordnumber = 0;
    uint charnumber = 0;


    //initializing character map for tokenizer
    QHash<char, bool> stopMap;
    char stopString[] = (" \r\n\t.;,:?!+=|()*[]\\|/<>&^%$#@`~\"'");
    for(uint i=0; i<strlen(stopString); i++)
    {
        stopMap[stopString[i]] = true;
    }

    //reading file in buffered way
    QFile file;
    if (argc>1) {
        file.setFileName(argv[1]);
    }
    else {
        cerr << "No input file specified" << endl;
        exit(1);
    }
    if (file.open(QFile::ReadOnly))
    {
        cout << "File exists !" << endl;
        QByteArray                 word;
        QByteArray                 qbuffer;
        QHash<QByteArray, quint64> map_words;
        char                       buf[1024];
        while (file.read(buf, 1024))
        {
            qbuffer = QByteArray(buf);

            foreach(char b, qbuffer)
            {
                charnumber++;
                if (!stopMap[b]==true)
                {
                    word.append(b);
                    //cout << "\n\nword  = " << word.data() << endl;
                }
                else
                {
                    if (word.length() > 0)
                    {

                        storeNgram1(word);
                        storeNgram2(word);
                        storeNgram3(word);
                        storeNgram4(word);

                        wordnumber++;
                        word.clear();
                    }
                    if (b == '\n')
                    {
                        linenumber++;
                    }
                }
            }
            word.clear();
        }
        file.close();
        cout << "---" << endl;
        cout << "map_words size = " << m4.size() << endl;
        cout << "Linenumber = " << linenumber << endl;
        cout << "wordnumber = " << wordnumber << endl;
        cout << "charnumber = " << charnumber << endl;
        cout << "---" << endl;
    }
    else
    {
        cout << "File missing !" << endl;
    }
    return 0;
}
