import os
import sys
import tensorflow as tf
from keras.models import load_model
import test_dataset.data_info as di

key = sys.argv[1]
test_model = sys.argv[2]
batch_size = 32

def pack_features_vector(features, labels):
  features = tf.stack(list(features.values()), axis=1)
  return features, labels

if key == "iris":
  test_data = 'test_dataset/iris_test.csv'
  label_name= di.iris_label_name
  label_names = di.iris_labels
  column_names = di.iris_columns
elif key == "cancer":
  test_data = 'test_dataset/cancer_test.csv'
  label_name= di.iris_label_name
  label_names = di.iris_labels
  column_names = di.iris_columns
elif key == "wine":
  test_data = 'test_dataset/wine_test.csv'
  label_name= di.iris_label_name
  label_names = di.iris_labels
  column_names = di.iris_columns
else:
  sys.exit("data not exist")     

test_dataset_fp = tf.keras.utils.get_file(fname=os.path.basename(test_data),origin=test_data)

test_dataset = tf.data.experimental.make_csv_dataset(test_dataset_fp,batch_size,column_names=column_names,label_name=label_name,num_epochs=1,shuffle=False)

test_dataset = test_dataset.map(pack_features_vector)
test_accuracy = tf.keras.metrics.Accuracy()

# model Acurracy
uploaded_model = load_model(test_model)

for (x, y) in test_dataset:
  logits = uploaded_model(x)
  prediction = tf.argmax(logits, axis=1, output_type=tf.int32)
  test_accuracy(prediction, y)

score = float(format(test_accuracy.result())) * 100
print(score)
