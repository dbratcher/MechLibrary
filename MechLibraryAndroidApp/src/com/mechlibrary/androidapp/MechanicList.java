package com.mechlibrary.androidapp;

import com.koushikdutta.async.future.FutureCallback;
import com.koushikdutta.async.future.Future;
import com.koushikdutta.ion.Ion;
import com.google.gson.JsonObject;
import com.google.gson.JsonArray;
import android.view.View.OnTouchListener;
import android.view.MotionEvent;
import android.content.Intent;
import android.os.Bundle;
import android.support.v4.app.Fragment;
import android.util.Log;
import android.view.LayoutInflater;
import android.view.View;
import android.view.ViewGroup;
import android.widget.ArrayAdapter;
import android.widget.Button;
import android.widget.ImageView;
import android.widget.ListView;
import android.widget.TextView;
import android.widget.Toast;

/**
 * A placeholder fragment containing a simple view.
 */
public class MechanicList extends Fragment {    
	
	ArrayAdapter<JsonObject> mechanicAdapter;
    Future<JsonArray> loading;
    String server = "http://mech-library.appspot.com";
    //String server = "http://10.0.2.2:8080";
    


    @Override
    public View onCreateView(LayoutInflater inflater, ViewGroup container,
            Bundle savedInstanceState) {
		View rootView = inflater.inflate(R.layout.mechanic_list, container, false);
		
		Button reloadButton = (Button) rootView.findViewById(R.id.reload);
		reloadButton.setOnClickListener(new View.OnClickListener() {
            public void onClick(View v) {
            	mechanicAdapter.clear();
            	load();
            }
        });
		mechanicAdapter = new ArrayAdapter<JsonObject>(getActivity(), 0) {
            @Override
            public View getView(int position, View convertView, ViewGroup parent) {
                if (convertView == null)
                    convertView = getActivity().getLayoutInflater().inflate(R.layout.mechanic, null);
                final JsonObject mechanic = getItem(position);
                Log.d("json object", mechanic.toString());

                // set the profile photo using Ion
                final String imageUrl = server + mechanic.get("ScreenshotURL").getAsString();

                ImageView imageView = (ImageView)convertView.findViewById(R.id.image);
                imageView.setOnTouchListener(new OnTouchListener() {
                    @Override
                    public boolean onTouch(View arg0, MotionEvent arg1) {
                        switch (arg1.getAction()) {
                        case MotionEvent.ACTION_DOWN: {
                        	Intent intent = new Intent(getActivity(), FullMechanicActivity.class);
                        	intent.putExtra("image", imageUrl);
                        	intent.putExtra("title", mechanic.get("Title").getAsString());
                        	intent.putExtra("description", mechanic.get("Description").getAsString());
                        	intent.putExtra("votes", mechanic.get("Votes").getAsString());
                    	    startActivity(intent);
                            break;
                        }
                        }
                        return true;
                    }
                });

                // Use Ion's builder set the google_image on an ImageView from a URL

                // start with the ImageView
                Ion.with(imageView)
                // load the url
                .load(imageUrl);

                // and finally, set the name and text
                TextView handle = (TextView)convertView.findViewById(R.id.title);
                handle.setText(mechanic.get("Title").getAsString());

                TextView text = (TextView)convertView.findViewById(R.id.description);
                text.setText(mechanic.get("Description").getAsString());
                return convertView;
            }
        };

        // basic setup of the ListView and adapter
        ListView listView = (ListView)rootView.findViewById(R.id.list);
        listView.setAdapter(mechanicAdapter);
        load();
		
        return rootView;
    }

    private void load() {
        // don't attempt to load more if a load is already in progress
        if (loading != null && !loading.isDone() && !loading.isCancelled())
            return;

        String url = server + "/android/mechanics/list/votes";

        Ion.getDefault(getActivity().getApplicationContext()).configure().setLogging("adapter", Log.DEBUG);
		loading = Ion.with(this)
        .load(url)
        .asJsonArray()
        .setCallback(new FutureCallback<JsonArray>() {
            @Override
            public void onCompleted(Exception e, JsonArray result) {
                // this is called back onto the ui thread, no Activity.runOnUiThread or Handler.post necessary.
                if (e != null) {
                    Toast.makeText(getActivity().getApplicationContext(), "Error loading mechanics.", Toast.LENGTH_LONG).show();
                    return;
                }
                for (int i = 0; i < result.size(); i++) {
                	mechanicAdapter.add(result.get(i).getAsJsonObject());
                }
            }
        });
    }
}
