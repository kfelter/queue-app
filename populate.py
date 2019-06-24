#/bin/python3
import sys, requests, time, random, datetime, json

courses = ['ECE', 'PHYS', 'COM', 'MATH', 'CPE', 'CS']
projects = ['Project', 'Lab', 'Essay', 'Assignment']
numbers = [101,100,201,200,301,300,401,400,501,500]
names = []
with open('static/names', 'r') as f:
    names = [line.strip() for line in f]

def queueRandom(host):
    firstName, lastName = random.choices(names, k=2)
    course = random.choice(courses) + '-' + str(random.choice(numbers))
    project = random.choice(projects) + '-' +  str(random.randint(0,16))
    email = '{}{}{}@vti.edc'.format(firstName, lastName, random.choice(range(80,99)))
    d = {
        "Name": firstName + ' ' + lastName,
        "Course": course,
        "Project": project,
        "Email": email
    }

    rd = requests.post(host+"/q/enqueue", json=d)
    print('enqueue', rd.json())

def dequeueOne(host):
    d = requests.get(host+"/q/dequeue")
    print('dequeue', d.json())
    if d.json() != None:
        return 1
    else:
        return 0

def populate(host):
    print("Populating:", host)
    while True:
        n = datetime.datetime.now()
        if n.hour > 7 and n.hour < 17 :
            numTAs = random.randint(1,2)
            numQueue = random.randint(0,3)

            for _ in range(numTAs):
                if dequeueOne(host) == 0:
                    for _ in range(3):
                        queueRandom(host)
            
            for _ in range(numQueue):
                queueRandom(host)


            time.sleep(5*60)
        else:
            print('Not in time frame')
            time.sleep(20*60)
    return

if __name__ == '__main__':

    host = "http://127.0.0.1:8080"
    if len(sys.argv) > 1:
        host = sys.argv[1]

    while dequeueOne(host):
        print('dequeue all')

    for _ in range(10):
        queueRandom(host)


    populate(host)