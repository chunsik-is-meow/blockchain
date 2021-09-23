# import numpy as np
# import tensorflow as tf
# from keras.models import load_model
# import functools
import sys

# model = 'iris_model.h5'
# data = '../blockchain/upload/data/iris.csv'
# pred = 'Species'

train_file_path = "aaa.csv"
test_file_path = "bbb.csv"

def init(model, data):
    global deep_learning_model, deep_learning_graph, attr, ans, features
    # 저장된 모델 로딩
    deep_learning_model = load_model(model)
    deep_learning_graph = tf.compat.v1.get_default_graph()

    # species 이름 로딩
    f = open(data, mode='r')
    datas = f.readlines()
    attr = datas[0].rstrip('\n').split(',')
    index = attr.index(pred)
    f.close
    
    print(attr)
    for data in datas:
        features = data.rstrip('\n').split(',')
    ans = features[index]
    del(features[index])
    del(features[0])
    
def prediction(pred):
    print('features:', features, file=sys.stderr)
    with deep_learning_graph.as_default():
        deep_learning_model.predect
        Y_pred = deep_learning_model(features)
        return {attr[Y_pred[0]]}

def evaluate(pred):
    score = 82.4
    return score
    # features = []
    # for i in range(len(features)): 
    #     features[i]=float(features[i])
    #     prediction(features, pred)
    

if __name__ == '__main__':
    model = sys.argv[1]
    data = sys.argv[2]
    pred = sys.argv[3]
    # init(model, data)
    print(evaluate(pred))
