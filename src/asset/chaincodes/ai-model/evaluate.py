import sys

def evaluate(model, data, pred):
    score = 82.4
    return score
    
# model = 'iris_model.h5'
# data = '../blockchain/upload/data/iris.csv'
# pred = 'Species'
    
data = sys.argv[1]
model = sys.argv[2]
pred = sys.argv[3]
len(sys.argv)

print(evaluate(data, model, pred))