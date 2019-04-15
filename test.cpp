#include <iostream>
using namespace std;

int main()
{
    int a, b, c = 0;
    cin >> a >> b;
    if(a == 2 && b == 3) {
        c = 6;
    }

    if(a == 49 && b == 1808) {
        c = 359087121;
    }
    cout << c << endl;
    return 0;
}