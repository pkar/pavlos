import sys
import os
import StringIO
import requests
import json
import random
from pymongo import MongoClient
from bson.objectid import ObjectId

try:
  from fabric.api import env, run, sudo, put, local, settings
  from fabric.colors import green, red, yellow
except ImportError as e:
  print "Install fabric first: sudo easy_install fabric"
  print "Run with: fab command"
  print "Unexpected error:", sys.exc_info()[0]
  raise

env.app           = 'pavlos'
env.tests         = [
  'models',
  'api',
]
env.apps          = {
  'src/main/pavlos.go': 'pavlos'
}
env.branch        = 'master'
env.git_repo      = 'git@github.com:%(app)s.git' % env
env.forward_agent = True
env.user          = 'admin'
env.key_filename  = '~/.somekeyfile.pem'

def staging():
  """
  env.hosts = ['']
  """
  env.stage = 'staging'
  env.hosts = 'pavl.com'
  env.branch = 'master'

def install():
  """

  """
  sudo('apt-key adv --keyserver keyserver.ubuntu.com --recv 7F0CEB10')
  sudo("echo 'deb http://downloads-distro.mongodb.org/repo/debian-sysvinit dist 10gen' | sudo tee /etc/apt/sources.list.d/mongodb.list")
  sudo('apt-get -y update')
  sudo('apt-get -y install gcc')
  sudo('apt-get install python-dev')
  sudo('apt-get install -y python-setuptools')
  sudo('apt-get -y install mongodb-10gen')
  sudo('apt-get install -y mongodb-10gen=2.4.6')
  #sudo('echo "mongodb-10gen hold" | sudo dpkg --set-selections')
  sudo('/etc/init.d/mongodb start')

  sudo('easy_install fabric')
  sudo('easy_install requests')
  sudo('easy_install pymongo')
  sudo('apt-get -y install supervisor')
  sudo('mkdir -p /var/log/supervisor')

  sudo('mkdir -p /var/apps/logs')
  sudo('mkdir -p /var/apps/static')
  sudo('chown %(user)s:%(user)s /var/apps' % env)

def deploy():
  """

  """
  build()
  sudo('chown %(user)s:%(user)s /var/apps' % env)
  run('mkdir -p /var/apps/logs')
  run('mkdir -p /var/apps/static')
  put('./config/supervisord.conf', '/tmp/supervisord.conf')
  sudo('mv /tmp/supervisord.conf /etc/supervisord.conf')
  put('./fabfile.py', '/var/apps/fabfile.py')
  put('./bin/pavlos', '/var/apps/pavlos_tmp')
  sudo('mv /var/apps/pavlos_tmp /var/apps/pavlos')
  run('chmod +x /var/apps/pavlos')
  sudo('rm -rf /var/apps/static')
  run('cd /var/apps; ./pavlos -load=true')
  run('sudo supervisorctl reload')
  run('sudo supervisorctl restart all')

def build():
  """
  Build binaries
  """
  packages()
  for path, exe in env.apps.iteritems():
    local('GOPATH=`pwd` GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o bin/%s %s' % (exe, path))
    local('echo 2013-09-29_12-14-13 > bin/BUILD')
    local('GOPATH=`pwd` echo go version >> bin/BUILD')

def lint():
  """
  Run golint on project, to install go get github.com/golang/lint/golint
  """
  local('golint src/main')

def test(x=None):
  """
  Run unit and integration tests
  """
  def p(result, name):
    if result.failed:
      print red("FAILED: " + name)
      print red('-'*80)
    else:
      print green("OK: " + name)
      print green('-'*80)
    print

  with settings(warn_only=True):
    if x:
      result = local('GOPATH=`pwd` go test ' + x)
      p(result, x)
      #result = local('GOPATH=`pwd` go vet ' + x)
      #p(result, x)
    else:
      for t in env.tests:
        result = local('GOPATH=`pwd` go test ' + t)
        p(result, t)
        #result = local('GOPATH=`pwd` go vet ' + t)
        #p(result, t)

def bench():
  """
  Run tests along with any benchmark tests defined.
  """
  for t in env.tests:
    result = local('GOPATH=`pwd` go test ' + t + ' -bench=.*')

def packages():
  """
  Extract all external packages to src/
  """
  with settings(warn_only=True):
    local('cd src; tar -xzvf github.com.tar.gz')
  with settings(warn_only=True):
    local('cd src; tar -xzvf code.google.com.tar.gz')
  with settings(warn_only=True):
    local('cd src; tar -xzvf launchpad.net.tar.gz')
  with settings(warn_only=True):
    local('cd src; tar -xzvf labix.org.tar.gz')

def packages_zip():
  """
  Compress all external packages
  """
  with settings(warn_only=True):
    local('cd src; tar -czvf github.com.tar.gz github.com')
  with settings(warn_only=True):
    local('cd src; tar -czvf code.google.com.tar.gz code.google.com')
  with settings(warn_only=True):
    local('cd src; tar -czvf launchpad.net.tar.gz launchpad.net')
  with settings(warn_only=True):
    local('cd src; tar -czvf labix.org.tar.gz labix.org')

def train_user_remote(name='pavlos', keywords=['greece', 'athens', 'travel']):
  """
  """
  if isinstance(keywords, list):
    keywords = ';'.join(keywords)
  run('cd /var/apps; fab train_user:%s,keywords="%s",port=80' % (name, keywords))


def train_user(name='pavlos', keywords=['greece', 'athens', 'travel'], port='8000'):
  """
  train a user to like keywords
  """
  if not isinstance(keywords, list):
    keywords = keywords.split(';')
  print name, keywords

  while True:
    print 'Getting recommendation /relevant/%s' % name
    r = requests.get('http://localhost:%s/relevant/%s' % (port, name))
    try:
      j = r.json()
    except:
      continue
    try:
      id = j['ID']
    except KeyError:
      continue

    like = random.randint(-1, 0)
    tlower = r.text.lower()
    if any(x in tlower for x in keywords):
      like = 1
      print 'Training %s to like /collect/%s' % (name, r.text)
    else:
      print 'Training %s to hate everything /collect/%s' % (name, name)

    r = requests.post('http://localhost:%s/collect/%s' % (port, name), data=json.dumps({"id": id, "like": like}))

def train_test_users():
  """
  train several users
  """
  jobs = [
    'fab train_user:pavlos,keywords="greece;athens;travel"',
    'fab train_user:tom,keywords="sports;football;baseball;hockey;travel;cuba;china;russia;syria"',
    'fab train_user:barry,keywords="jesus;god;satan;religion"',
    'fab train_user:bob,keywords="obama"',
  ]

  for job in jobs:
    local(job)

def dump_user_remote(name='pavlos'):
  """
  show data on remote user
  """
  run('cd /var/apps; fab dump_user:%s' % name)

def dump_user(name='pavlos'):
  """
  dump user information
  """
  client = MongoClient()
  db = client['pavlos']
  users = db['Users']
  items = db['Items']
  categories = db['Categories']

  user = users.find_one({'name': name})
  print 80 * '*'
  print user['name']
  print 80 * '*'
  for id, meta in user['categories'].iteritems():
    category = categories.find_one({'_id': ObjectId(id)})
    if category:
      print category['displayname'], 'Weight:', meta['weight'], 'Score:', meta['score'], 'Count:', meta['count']

  print 2 * '\n'
  for id, meta in user['sims'].iteritems():
    other = users.find_one({'_id': ObjectId(id)})
    print other['name'], 'Euclidean:', meta['euclidean'], 'Pearson:', meta['pearson']


def load_remote_data():
  """
  show data on remote user
  """
  run('cd /var/apps; ./pavlos -load=true -logtostderr=true')

